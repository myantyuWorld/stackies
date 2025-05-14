package presenter_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"stackies/backend/presenter"
	"stackies/backend/usecase"
	mock_usecase "stackies/backend/usecase/mock"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestExperienceHandler_Create(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(mock *mock_usecase.MockExperienceUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系: 体験を作成できる",
			requestBody: `{"title":"テスト体験"}`,
			setupMock: func(mock *mock_usecase.MockExperienceUsecase) {
				mock.EXPECT().Create("テスト体験").Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"title":"テスト体験"}`,
		},
		{
			name:           "異常系: リクエストボディが不正",
			requestBody:    `{"title":123}`,
			setupMock:      func(mock *mock_usecase.MockExperienceUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoのインスタンスを作成
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/experiences", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// モックの設定
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUsecase := mock_usecase.NewMockExperienceUsecase(ctrl)
			tt.setupMock(mockUsecase)

			// ハンドラーの作成
			handler := presenter.NewExperienceHandler(mockUsecase)

			// テスト対象の実行
			err := handler.Create(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestExperienceHandler_GetAll(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		setupMock      func(mock *mock_usecase.MockExperienceUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系: 体験一覧を取得できる",
			setupMock: func(mock *mock_usecase.MockExperienceUsecase) {
				mock.EXPECT().GetAll().Return([]usecase.ExperienceDto{
					{ID: 1, Title: "体験1"},
					{ID: 2, Title: "体験2"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"title":"体験1"},{"id":2,"title":"体験2"}]`,
		},
		{
			name: "異常系: 体験一覧の取得に失敗",
			setupMock: func(mock *mock_usecase.MockExperienceUsecase) {
				mock.EXPECT().GetAll().Return(nil, errors.New("データベースエラー"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"データベースエラー"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoのインスタンスを作成
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/experiences", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// モックの設定
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUsecase := mock_usecase.NewMockExperienceUsecase(ctrl)
			tt.setupMock(mockUsecase)

			// ハンドラーの作成
			handler := presenter.NewExperienceHandler(mockUsecase)

			// テスト対象の実行
			err := handler.GetAll(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}
