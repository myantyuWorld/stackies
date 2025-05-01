package main

import (
	"net/http"

	"stackies/backend/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// データベース接続の初期化
	dbConfig := config.NewDBConfig()
	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	// ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Stackies API",
		})
	})

	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
