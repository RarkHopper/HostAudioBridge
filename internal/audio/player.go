package audio

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FilePath string

func (f FilePath) Exists() bool {
	_, err := os.Stat(string(f))
	return !os.IsNotExist(err)
}

type DirPath string

func (d DirPath) Exists() bool {
	info, err := os.Stat(string(d))
	return err == nil && info.IsDir()
}

var (
	ErrInvalidAudioName = errors.New("不正な音声名")
	ErrInvalidVolume    = errors.New("音量は0.0〜1.0の範囲で指定してください")
)

// Audio は音声名を表す値オブジェクト
type Audio string

var validAudioName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func NewAudio(s string) (Audio, error) {
	if !validAudioName.MatchString(s) {
		return "", ErrInvalidAudioName
	}
	return Audio(s), nil
}

// Volume は音量を表す値オブジェクト
type Volume float64

func NewVolume(v float64) (Volume, error) {
	if v < 0 || v > 1 {
		return 0, ErrInvalidVolume
	}
	return Volume(v), nil
}

type Player interface {
	Play(ctx context.Context, audio Audio, vol Volume) error
	List() []Audio
}

// ScanAudioDir はディレクトリから利用可能な音声を取得する
func ScanAudioDir(dir DirPath) []Audio {
	pattern := filepath.Join(string(dir), "*.wav")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return []Audio{}
	}

	list := make([]Audio, 0, len(files))
	for _, f := range files {
		name := strings.TrimSuffix(filepath.Base(f), ".wav")
		if a, err := NewAudio(name); err == nil {
			list = append(list, a)
		}
	}
	return list
}
