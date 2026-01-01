//go:build darwin

package audio

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"slices"

	"github.com/rarkhopper/host-audio-bridge/internal/fs"
)

var (
	// ErrFileNotFound は音声ファイルが見つからないことを示す
	ErrFileNotFound = errors.New("音声ファイルが見つかりません")
	// ErrAudioDirNotFound は音声ディレクトリが見つからないことを示す
	ErrAudioDirNotFound = errors.New("音声ディレクトリが見つかりません")
)

type darwinPlayer struct {
	audioDir       fs.DirPath
	availableAudio []Audio
}

// NewPlayer は新しいdarwinプレイヤーを生成する
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
		if err := exec.Command("afplay", "-v", fmt.Sprintf("%.2f", vol), string(fp)).Run(); err != nil {
			log.Printf("音声再生エラー: %v", err)
		}
	}()
	return nil
}

func (p *darwinPlayer) List() []Audio {
	return slices.Clone(p.availableAudio)
}
