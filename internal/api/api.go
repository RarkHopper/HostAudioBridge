package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rarkhopper/host-audio-bridge/internal/audio"
)

const defaultVolume = 1.0

type PlayRequest struct {
	Audio  string   `json:"audio"`
	Volume *float64 `json:"volume,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type AudioListResponse struct {
	Audio []audio.Audio `json:"audio"`
}

func RegisterRoutes(e *echo.Echo, p audio.Player) {
	e.GET("/health", handleHealth())
	e.POST("/play", handlePlay(p))
	e.GET("/audio", handleAudioList(p))
}

func handleHealth() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}

func handlePlay(p audio.Player) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req PlayRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "リクエストの形式が不正です",
			})
		}

		a, err := audio.NewAudio(req.Audio)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
		}

		v := defaultVolume
		if req.Volume != nil {
			v = *req.Volume
		}
		vol, err := audio.NewVolume(v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
		}

		if err := p.Play(c.Request().Context(), a, vol); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "再生に失敗しました: " + err.Error(),
			})
		}

		return c.NoContent(http.StatusOK)
	}
}

func handleAudioList(p audio.Player) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, AudioListResponse{
			Audio: p.List(),
		})
	}
}
