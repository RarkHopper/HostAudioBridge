//go:build darwin

package audio

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
)

var (
	ErrFileNotFound     = errors.New("音声ファイルが見つかりません")
	ErrAudioDirNotFound = errors.New("音声ディレクトリが見つかりません")
)

type darwinPlayer struct {
	audioDir       DirPath
	availableAudio []Audio
}

func NewPlayer(audioDir DirPath, availableAudio []Audio) (Player, error) {
	if !audioDir.Exists() {
		return nil, ErrAudioDirNotFound
	}
	return &darwinPlayer{
		audioDir:       audioDir,
		availableAudio: availableAudio,
	}, nil
}

func (p *darwinPlayer) Play(ctx context.Context, audio Audio, vol Volume) error {
	fp := FilePath(filepath.Join(string(p.audioDir), string(audio)+".wav"))
	if !fp.Exists() {
		return ErrFileNotFound
	}

	//#nosec G204 -- audio は NewAudio でバリデーション済み
	if err := exec.CommandContext(ctx, "afplay", "-v", fmt.Sprintf("%.2f", vol), string(fp)).Run(); err != nil {
		return fmt.Errorf("音声ファイルの再生に失敗: %w", err)
	}
	return nil
}

func (p *darwinPlayer) List() []Audio {
	return slices.Clone(p.availableAudio)
}
