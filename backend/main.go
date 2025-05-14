package main

import (
	"fmt"
	"net/http"

	"stackies/backend/config"
	"stackies/backend/infra/repository"
	"stackies/backend/presenter"
	"stackies/backend/usecase"

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

	e.POST("/experiences", experienceHandler.Create)
	// http://localhost:28080/experiences
	e.GET("/experiences", experienceHandler.GetAll)

	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
