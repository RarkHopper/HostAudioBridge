package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rarkhopper/host-audio-bridge/internal/audio"
)

type playRequest struct {
	Audio  string   `json:"audio"`
	Volume *float64 `json:"volume"`
}

const defaultVolume = 1.0

type playResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type audioListResponse struct {
	Status string        `json:"status"`
	Audio  []audio.Audio `json:"audio"`
}

func RegisterRoutes(e *echo.Echo, p audio.Player) {
	e.POST("/play", handlePlay(p))
	e.GET("/audio", handleAudioList(p))
}

func handlePlay(p audio.Player) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req playRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, playResponse{
				Status:  "error",
				Message: "リクエストの形式が不正です",
			})
		}

		a, err := audio.NewAudio(req.Audio)
		if err != nil {
			return c.JSON(http.StatusBadRequest, playResponse{
				Status:  "error",
				Message: err.Error(),
			})
		}

		v := defaultVolume
		if req.Volume != nil {
			v = *req.Volume
		}
		vol, err := audio.NewVolume(v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, playResponse{
				Status:  "error",
				Message: err.Error(),
			})
		}

		if err := p.Play(c.Request().Context(), a, vol); err != nil {
			return c.JSON(http.StatusInternalServerError, playResponse{
				Status:  "error",
				Message: "再生に失敗しました: " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, playResponse{Status: "ok"})
	}
}

func handleAudioList(p audio.Player) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, audioListResponse{
			Status: "ok",
			Audio:  p.List(),
		})
	}
}
