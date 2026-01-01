//go:build darwin

package audio

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"

	"github.com/rarkhopper/host-audio-bridge/internal/fs"
)

var (
	ErrFileNotFound     = errors.New("音声ファイルが見つかりません")
	ErrAudioDirNotFound = errors.New("音声ディレクトリが見つかりません")
)

type darwinPlayer struct {
	audioDir       fs.DirPath
	availableAudio []Audio
}

func NewPlayer(audioDir fs.DirPath, availableAudio []Audio) (Player, error) {
	if !audioDir.Exists() {
		return nil, ErrAudioDirNotFound
	}
	return &darwinPlayer{
		audioDir:       audioDir,
		availableAudio: availableAudio,
	}, nil
}

func (p *darwinPlayer) Play(_ context.Context, audio Audio, vol Volume) error {
	fp := p.audioDir.Join(string(audio) + ".wav")
	if !fp.Exists() {
		return ErrFileNotFound
	}

	go func() {
		//#nosec G204 -- audio は NewAudio でバリデーション済み
		_ = exec.Command("afplay", "-v", fmt.Sprintf("%.2f", vol), string(fp)).Run()
	}()
	return nil
}

func (p *darwinPlayer) List() []Audio {
	return slices.Clone(p.availableAudio)
}
