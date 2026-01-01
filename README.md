# HostAudioBridge

リモート環境からホストマシンで音声を再生するためのブリッジツール。

## 概要

オーディオデバイスにアクセスできない環境（DevContainer、Docker、SSH先、WSLなど）から、HTTP API経由でホスト側に音声再生を委譲する。

```
┌─────────────────┐     HTTP      ┌─────────────────┐
│  Remote Env     │ ──────────▶  │   Host Machine  │
│                 │              │                 │
│  hab-cli        │   POST /play │  server         │
│                 │              │  └─▶ afplay     │
└─────────────────┘              └─────────────────┘
```

## なぜ作ったか

PulseAudio/PipeWireのソケット転送は設定が複雑で、WSLgはWindows限定。HTTPで音声名を送るだけのシンプルなツールが欲しかった。

## 動作環境

- **ホスト**: macOS（afplayコマンドを使用）
- **Go**: 1.25以上

> [!WARNING]
> 個人使用のローカルツールのため、Windows/Linux対応は未定

## セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/RarkHopper/HostAudioBridge.git
cd HostAudioBridge/

# 環境変数ファイルを作成
make env-init

# ビルド
make build
```

## 使い方

### サーバー（ホスト側）

```bash
make up
```

デフォルトでポート47800で起動。

### CLI（リモート側）

```bash
make cli
```

対話形式で音声を選択・再生できる。

## 環境変数

| 変数名 | 説明 |
|--------|------|
| `HAB_PORT` | サーバーのポート番号 |
| `HAB_SERVER_URL` | CLIが接続するサーバーURL |

`.env.example`を参考に`.env`を作成：

```bash
HAB_PORT=47800
HAB_SERVER_URL=http://host.docker.internal:47800
```

## 音声ファイル

[audio/README.md](audio/README.md) を参照。

## 開発

```bash
make test    # テスト
make check   # Lint
make format  # フォーマット
```
