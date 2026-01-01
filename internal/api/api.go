// api パッケージはHTTPハンドラを提供する
package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes はルートをEchoインスタンスに登録する
func RegisterRoutes(e *echo.Echo) {
	e.POST("/play", handlePlay)
	e.GET("/sounds", handleSounds)
}

// handlePlay は音声再生リクエストを処理する
func handlePlay(c echo.Context) error {
	// TODO: 実装
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// handleSounds はプリセット音声の一覧を返す
func handleSounds(c echo.Context) error {
	// TODO: sounds ディレクトリから動的に読み込む
	sounds := []string{"notification", "error", "success", "bell"}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"sounds": sounds,
	})
}
