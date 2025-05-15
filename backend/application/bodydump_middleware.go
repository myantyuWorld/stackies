package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// BodyDumpConfig は BodyDumpMiddleware の設定です
type BodyDumpConfig struct {
	// Skipper は特定のリクエストでミドルウェアをスキップするための関数です
	Skipper middleware.Skipper

	// MaxBodySize はログに出力するボディの最大サイズです（バイト単位）
	// 0は制限なしを意味します
	MaxBodySize int64

	// ExcludeURLs はボディをログに出力しないURLパターンのリストです
	ExcludeURLs []string

	// MaskFields は機密情報をマスクするフィールド名のリストです
	MaskFields []string

	// Logger はカスタムロガーです
	Logger echo.Logger
}

// DefaultBodyDumpConfig はデフォルトの設定です
var DefaultBodyDumpConfig = BodyDumpConfig{
	Skipper:     middleware.DefaultSkipper,
	MaxBodySize: 1024 * 10, // 10KB
	ExcludeURLs: []string{},
	MaskFields:  []string{"password", "token", "secret", "authorization", "api_key"},
	Logger:      nil,
}

// BodyDump はリクエストとレスポンスのボディをログに出力するミドルウェアです
func BodyDump(next echo.HandlerFunc) echo.HandlerFunc {
	return BodyDumpWithConfig(DefaultBodyDumpConfig)(next)
}

// BodyDumpWithConfig はカスタム設定でボディをログに出力するミドルウェアです
func BodyDumpWithConfig(config BodyDumpConfig) echo.MiddlewareFunc {
	// デフォルト値の設定
	if config.Skipper == nil {
		config.Skipper = DefaultBodyDumpConfig.Skipper
	}
	if config.Logger == nil {
		config.Logger = log.New("body-dump")
	}
	if config.MaxBodySize == 0 {
		config.MaxBodySize = DefaultBodyDumpConfig.MaxBodySize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// URLがExcludeURLsに含まれているかチェック
			for _, excludeURL := range config.ExcludeURLs {
				if strings.Contains(c.Request().URL.Path, excludeURL) {
					return next(c)
				}
			}

			// リクエストボディの読み取りと復元
			var reqBody []byte
			if c.Request().Body != nil {
				// リクエストボディの読み取り
				limitedReader := io.LimitReader(c.Request().Body, config.MaxBodySize)
				reqBody, _ = io.ReadAll(limitedReader)

				// リクエストボディの復元
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			// レスポンスボディをキャプチャするためのレスポンスライター
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{
				ResponseWriter: c.Response().Writer,
				Writer:         mw,
			}
			c.Response().Writer = writer

			// リクエスト開始時間
			start := time.Now()

			// 次のハンドラを実行
			err := next(c)

			// レスポンス時間
			duration := time.Since(start)

			// リクエスト情報
			req := c.Request()
			res := c.Response()

			// リクエストとレスポンスのボディをマスク処理
			maskedReqBody := maskSensitiveData(reqBody, config.MaskFields)
			maskedResBody := maskSensitiveData(resBody.Bytes(), config.MaskFields)

			// ログ出力
			logData := map[string]interface{}{
				"time":             time.Now().Format(time.RFC3339),
				"remote_ip":        c.RealIP(),
				"host":             req.Host,
				"method":           req.Method,
				"uri":              req.RequestURI,
				"user_agent":       req.UserAgent(),
				"status":           res.Status,
				"duration":         duration.String(),
				"duration_ms":      float64(duration.Nanoseconds()) / 1e6,
				"request_headers":  formatHeaders(req.Header),
				"response_headers": formatHeaders(res.Header()),
			}

			// Content-Typeがapplication/jsonの場合は整形して出力
			reqContentType := req.Header.Get("Content-Type")
			resContentType := res.Header().Get("Content-Type")

			if strings.Contains(reqContentType, "application/json") && len(maskedReqBody) > 0 {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, maskedReqBody, "", "  "); err == nil {
					logData["request_body"] = prettyJSON.String()
				} else {
					logData["request_body"] = string(maskedReqBody)
				}
			} else if len(maskedReqBody) > 0 {
				logData["request_body"] = string(maskedReqBody)
			}

			if strings.Contains(resContentType, "application/json") && len(maskedResBody) > 0 {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, maskedResBody, "", "  "); err == nil {
					logData["response_body"] = prettyJSON.String()
				} else {
					logData["response_body"] = string(maskedResBody)
				}
			} else if len(maskedResBody) > 0 {
				logData["response_body"] = string(maskedResBody)
			}

			// ステータスコードに応じてログレベルを変更
			if res.Status >= 500 {
				config.Logger.Errorf("%+v", logData)
			} else if res.Status >= 400 {
				config.Logger.Warnf("%+v", logData)
			} else {
				config.Logger.Infof("%+v", logData)
			}

			return err
		}
	}
}

// bodyDumpResponseWriter はレスポンスボディをキャプチャするためのラッパーです
type bodyDumpResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// maskSensitiveData は機密データをマスクする関数です
func maskSensitiveData(data []byte, maskFields []string) []byte {
	if len(data) == 0 {
		return data
	}

	// JSONデータかどうかをチェック
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		// JSONでない場合はそのまま返す
		return data
	}

	// 機密フィールドをマスク
	for _, field := range maskFields {
		maskJSONField(jsonData, field, "********")
	}

	// マスク済みのJSONを返す
	maskedData, err := json.Marshal(jsonData)
	if err != nil {
		return data // エラーが発生した場合は元のデータを返す
	}

	return maskedData
}

// maskJSONField はJSONオブジェクト内の特定のフィールドをマスクする関数です
func maskJSONField(data map[string]interface{}, field, mask string) {
	for key, value := range data {
		if strings.ToLower(key) == strings.ToLower(field) {
			data[key] = mask
			continue
		}

		// ネストされたオブジェクトを再帰的に処理
		switch v := value.(type) {
		case map[string]interface{}:
			maskJSONField(v, field, mask)
		case []interface{}:
			for _, item := range v {
				if mapItem, ok := item.(map[string]interface{}); ok {
					maskJSONField(mapItem, field, mask)
				}
			}
		}
	}
}

// formatHeaders はヘッダーを文字列に整形する関数です
func formatHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		// 認証ヘッダーはマスク
		if strings.ToLower(key) == "authorization" || strings.ToLower(key) == "cookie" {
			result[key] = "********"
		} else {
			result[key] = strings.Join(values, ", ")
		}
	}
	return result
}
