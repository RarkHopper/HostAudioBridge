package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rarkhopper/host-audio-bridge/internal/api"
	"github.com/rarkhopper/host-audio-bridge/internal/audio"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("HAB_PORT")
	if port == "" {
		log.Fatal("HAB_PORT 環境変数が設定されていません")
	}

	e := echo.New()
	e.HideBanner = true

	setupMiddleware(e)

	audioDir := audio.DirPath("audio")
	availableAudio := audio.ScanAudioDir(audioDir)
	player, err := audio.NewPlayer(audioDir, availableAudio)
	if err != nil {
		log.Fatalf("プレイヤー初期化エラー: %v", err)
	}
	api.RegisterRoutes(e, player)

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
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%s %s %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
}
