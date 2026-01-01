// Package audio は音声再生のドメインロジックを提供する
package audio

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/rarkhopper/host-audio-bridge/internal/fs"
)

var (
	// ErrInvalidAudioName は不正な音声名を示す
	ErrInvalidAudioName = errors.New("不正な音声名")
	// ErrInvalidVolume は不正な音量値を示す
	ErrInvalidVolume = errors.New("音量は0.0〜1.0の範囲で指定してください")
)

// Audio は音声名を表す値オブジェクト
type Audio string

var validAudioName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// NewAudio は音声名を検証してAudioを生成する
func NewAudio(s string) (Audio, error) {
	if !validAudioName.MatchString(s) {
		return "", ErrInvalidAudioName
	}
	return Audio(s), nil
}

// Volume は音量を表す値オブジェクト
type Volume float64

// NewVolume は音量値を検証してVolumeを生成する
func NewVolume(v float64) (Volume, error) {
	if v < 0 || v > 1 {
		return 0, ErrInvalidVolume
	}
	return Volume(v), nil
}

// Player は音声再生のインターフェース
type Player interface {
	Play(ctx context.Context, audio Audio, vol Volume) error
	List() []Audio
}

// ScanAudioDir はディレクトリから利用可能な音声を取得する
func ScanAudioDir(dir fs.DirPath) []Audio {
	files, err := dir.Glob("*.wav")
	if err != nil {
		return []Audio{}
	}

	list := make([]Audio, 0, len(files))
	for _, f := range files {
		name := strings.TrimSuffix(f.Base(), ".wav")
		if a, err := NewAudio(name); err == nil {
			list = append(list, a)
		}
	}
	return list
}
