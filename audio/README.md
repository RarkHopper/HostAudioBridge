# Audio Files

このディレクトリに再生する音声ファイルを配置してください。

## 要件

| 項目 | 要件 |
|------|------|
| 形式 | WAV (.wav) |
| ファイル名 | 英数字、ハイフン、アンダースコアのみ |
| 例 | `bell.wav`, `notify-01.wav`, `alert_sound.wav` |

## 無効なファイル名

以下のファイル名は無視されます:

- スペースを含む: `my sound.wav`
- 日本語: `通知音.wav`
- 特殊文字: `alert!.wav`

## 推奨ファイル

| 用途 | ファイル名例 |
|------|-------------|
| 通知 | `notification.wav` |
| 成功 | `success.wav` |
| エラー | `error.wav` |
| 警告 | `warning.wav` |
