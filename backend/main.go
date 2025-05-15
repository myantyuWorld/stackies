package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"stackies/backend/config"
	"stackies/backend/infra/repository"
	"stackies/backend/presenter"
	"stackies/backend/usecase"

	"github.com/MicahParks/keyfunc"
	"github.com/coreos/go-oidc"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURL  = os.Getenv("REDIRECT_URL")
	issuerURL    = os.Getenv("ISSUER_URL")
	tokenURL     = os.Getenv("TOKEN_URL")
	provider     *oidc.Provider
	oauth2Config oauth2.Config
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

const (
	jwksURL = "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_pQJjZxAvp/.well-known/jwks.json"
)

var jwks *keyfunc.JWKS

func main() {
	var err error
	jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "JWKs取得失敗: %v\n", err)
		os.Exit(1)
	}

	// Echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BodyDump(func(c echo.Context, req []byte, res []byte) {
		fmt.Println("Request Body:", string(req))
		fmt.Println("Response Body:", string(res))
	}))

	// データベース接続の初期化
	dbConfig := config.NewDBConfig()
	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}

	experienceRepository := repository.NewExperienceRepository(db)
	experienceUsecase := usecase.NewExperienceUsecase(experienceRepository)
	experienceHandler := presenter.NewExperienceHandler(experienceUsecase)

	// ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Stackies API",
		})
	})

	// コールバックエンドポイント
	e.GET("/callback", func(c echo.Context) error {
		// クエリパラメータからcodeやstateを取得
		code := c.QueryParam("code")
		errorMsg := c.QueryParam("error")

		if errorMsg != "" {
			return c.String(http.StatusBadRequest, fmt.Sprintf("認証エラー: %s", errorMsg))
		}

		if code == "" {
			return c.String(http.StatusBadRequest, "codeがありません")
		}

		spew.Dump(clientID)
		spew.Dump(clientSecret)
		spew.Dump(redirectURL)
		spew.Dump(issuerURL)
		spew.Dump(tokenURL)

		// トークンエンドポイントにPOST
		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("client_id", clientID)
		data.Set("code", code)
		data.Set("redirect_uri", redirectURL)

		req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
		if err != nil {
			return c.String(http.StatusInternalServerError, "リクエスト作成失敗")
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Basic認証（client_secretが必要な場合のみ）
		if clientSecret != "" {
			basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
			req.Header.Set("Authorization", "Basic "+basicAuth)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.String(http.StatusInternalServerError, "トークンリクエスト失敗")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("トークン取得失敗: %s", string(body)))
		}

		var tokenRes TokenResponse
		if err := json.Unmarshal(body, &tokenRes); err != nil {
			return c.String(http.StatusInternalServerError, "トークンレスポンスのパース失敗")
		}

		// 取得したトークンを表示
		return c.JSON(http.StatusOK, tokenRes)

	})

	e.POST("/experiences", experienceHandler.Create, JWTMiddleware)
	// http://localhost:28080/experiences
	e.GET("/experiences", experienceHandler.GetAll, JWTMiddleware)

	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}

// JWT認証ミドルウェア
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, "トークンがありません")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "トークンが無効です")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, "クレーム取得失敗")
		}

		// 必要ならaudやissの検証もここで追加
		// 例:
		// if claims["aud"] != "your-client-id" { ... }

		// claimsをコンテキストにセット
		c.Set("claims", claims)
		spew.Dump(claims)

		return next(c)
	}
}
