package audio

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewAudio はAudio値オブジェクトのバリデーションをテストする。
// 有効な音声名（英数字、ハイフン、アンダースコア）のみ許可し、
// コマンドインジェクションに繋がる特殊文字を拒否することを確認する。
func TestNewAudio(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "bell", false},
		{"valid with hyphen", "notify-sound", false},
		{"valid with underscore", "notify_sound", false},
		{"valid with numbers", "sound123", false},
		{"valid mixed", "My-Sound_01", false},
		{"empty string", "", true},
		{"with space", "bell sound", true},
		{"with dot", "bell.wav", true},
		{"with slash", "path/bell", true},
		{"with backslash", "path\\bell", true},
		{"japanese", "通知音", true},
		{"special chars", "bell!@#", true},
		{"semicolon injection", "bell;rm -rf", true},
		{"pipe injection", "bell|cat /etc/passwd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAudio(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAudio(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// TestNewVolume はVolume値オブジェクトのバリデーションをテストする。
// 0.0（ミュート）から1.0（最大）の範囲のみ許可することを確認する。
func TestNewVolume(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr bool
	}{
		{"zero (mute)", 0, false},
		{"mid volume", 0.5, false},
		{"max volume", 1, false},
		{"small positive", 0.01, false},
		{"near max", 0.99, false},
		{"negative", -0.1, true},
		{"over max", 1.1, true},
		{"large negative", -100, true},
		{"large positive", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVolume(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVolume(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// TestScanAudioDir はディレクトリスキャンのバリデーション統合をテストする。
// 有効な音声名を持つ.wavファイルのみ返し、無効なファイル名は除外することを確認する。
func TestScanAudioDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "audio_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// テスト用ファイルを作成
	validFiles := []string{"bell.wav", "notify.wav", "alert-01.wav"}
	invalidFiles := []string{"has space.wav", "日本語.wav", "no_extension"}

	for _, f := range append(validFiles, invalidFiles...) {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte{}, 0600); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// テスト実行
	result := ScanAudioDir(DirPath(tmpDir))

	// 有効なファイルのみ返されることを確認
	if len(result) != len(validFiles) {
		t.Errorf("ScanAudioDir() returned %d files, want %d", len(result), len(validFiles))
	}

	// 各有効ファイルが含まれていることを確認
	resultMap := make(map[Audio]bool)
	for _, a := range result {
		resultMap[a] = true
	}

	for _, f := range validFiles {
		name := f[:len(f)-4] // .wav を除去
		if !resultMap[Audio(name)] {
			t.Errorf("ScanAudioDir() missing expected audio: %s", name)
		}
	}
}

// TestScanAudioDir_EmptyDir は空ディレクトリに対して空スライスを返すことを確認する。
func TestScanAudioDir_EmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "audio_test_empty")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	result := ScanAudioDir(DirPath(tmpDir))

	if len(result) != 0 {
		t.Errorf("ScanAudioDir() on empty dir returned %d files, want 0", len(result))
	}
}

// TestScanAudioDir_NonExistent は存在しないディレクトリに対して空スライスを返すことを確認する。
func TestScanAudioDir_NonExistent(t *testing.T) {
	result := ScanAudioDir(DirPath("/nonexistent/path"))

	if len(result) != 0 {
		t.Errorf("ScanAudioDir() on nonexistent dir returned %d files, want 0", len(result))
	}
}
