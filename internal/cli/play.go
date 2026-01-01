package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/rarkhopper/host-audio-bridge/internal/audio"
)

func (c *CLI) handlePlay(ctx context.Context) {
	audioList, err := c.client.List(ctx)
	if err != nil {
		fmt.Printf("音声一覧の取得に失敗: %v\n", err)
		return
	}

	if len(audioList) == 0 {
		fmt.Println("利用可能な音声ファイルがありません")
		return
	}

	var selectedAudio audio.Audio
	var volumeStr string

	options := make([]huh.Option[audio.Audio], len(audioList))
	for i, a := range audioList {
		options[i] = huh.NewOption(string(a), a)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[audio.Audio]().
				Title("音声を選択").
				Options(options...).
				Value(&selectedAudio),
			huh.NewInput().
				Title("音量 (0.0-1.0, 空欄でデフォルト1.0)").
				Value(&volumeStr).
				Validate(validateVolume),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}

	var vol *audio.Volume
	if volumeStr != "" {
		// validateVolumeで検証済み
		v, _ := strconv.ParseFloat(volumeStr, 64) //nolint:errcheck
		parsed, _ := audio.NewVolume(v)           //nolint:errcheck
		vol = &parsed
	}

	if vol != nil {
		fmt.Printf("再生中: %s (音量: %.1f)\n", selectedAudio, *vol)
	} else {
		fmt.Printf("再生中: %s\n", selectedAudio)
	}

	if err := c.client.Play(ctx, selectedAudio, vol); err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}

	fmt.Println("完了")
}

func validateVolume(s string) error {
	if s == "" {
		return nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return audio.ErrInvalidVolume
	}

	_, err = audio.NewVolume(v)
	return err //nolint:wrapcheck // バリデーションメッセージとして直接表示
}
