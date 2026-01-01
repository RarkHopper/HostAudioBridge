package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rarkhopper/host-audio-bridge/internal/api"
	"github.com/rarkhopper/host-audio-bridge/internal/audio"
)

var (
	ErrUnexpectedStatus = errors.New("予期しないステータスコード")
	ErrServerStatus     = errors.New("サーバーがエラーステータスを返しました")
)

type Client interface {
	List(ctx context.Context) ([]audio.Audio, error)
	Play(ctx context.Context, a audio.Audio, vol *audio.Volume) error
	Health(ctx context.Context) error
}

type httpClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *httpClient) List(ctx context.Context) ([]audio.Audio, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/audio", nil)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("音声一覧の取得に失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}

	var result api.AudioListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードに失敗: %w", err)
	}

	return result.Audio, nil
}

func (c *httpClient) Play(ctx context.Context, a audio.Audio, vol *audio.Volume) error {
	reqBody := api.PlayRequest{
		Audio: string(a),
	}
	if vol != nil {
		v := float64(*vol)
		reqBody.Volume = &v
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("リクエストのマーシャルに失敗: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/play", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("音声の再生に失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp api.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Message != "" {
			return errors.New(errResp.Message)
		}
		return fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}

	return nil
}

func (c *httpClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("サーバーに接続できません: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrServerStatus, resp.StatusCode)
	}

	return nil
}
