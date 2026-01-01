// HostAudioBridge サーバーのエントリーポイント
package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rarkhopper/host-audio-bridge/internal/api"
)

func main() {
	port := os.Getenv("HAB_PORT")
	if port == "" {
		log.Fatal("HAB_PORT 環境変数が設定されていません")
	}

	e := echo.New()
	e.HideBanner = true

	setupMiddleware(e)
	api.RegisterRoutes(e)

	log.Printf("サーバーを起動します: :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}

func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%s %s %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
}
