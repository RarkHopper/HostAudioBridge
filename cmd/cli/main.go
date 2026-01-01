package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rarkhopper/host-audio-bridge/internal/cli"
)

const envServerURL = "HAB_SERVER_URL"

func main() {
	_ = godotenv.Load()

	serverURL := os.Getenv(envServerURL)
	if serverURL == "" {
		fmt.Printf("%s 環境変数が設定されていません\n", envServerURL)
		os.Exit(1)
	}

	client := cli.NewClient(serverURL)
	app := cli.New(client)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
