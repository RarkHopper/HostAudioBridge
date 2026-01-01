package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

var ErrServerConnection = errors.New("サーバーに接続できません")

type Feature struct {
	DisplayName string
	Execute     func(ctx context.Context)
	IsExit      bool
}

type CLI struct {
	client   Client
	features []Feature
}

func New(client Client) *CLI {
	c := &CLI{client: client}
	c.features = []Feature{
		{DisplayName: "音声を再生", Execute: c.handlePlay},
		{DisplayName: "音声一覧を表示", Execute: c.handleList},
		{DisplayName: "終了", IsExit: true},
	}
	return c
}

func (c *CLI) Run() error {
	fmt.Println("HostAudioBridge CLI")
	fmt.Println()

	ctx := context.Background()

	fmt.Print("サーバーに接続中... ")
	if err := c.client.Health(ctx); err != nil {
		fmt.Println("失敗")
		return fmt.Errorf("%w: %w", ErrServerConnection, err)
	}
	fmt.Println("接続完了")
	fmt.Println()

	for {
		if !c.runLoop(ctx) {
			break
		}
	}

	fmt.Println("終了します")
	return nil
}

func (c *CLI) runLoop(ctx context.Context) bool {
	options := make([]huh.Option[*Feature], len(c.features))
	for i := range c.features {
		options[i] = huh.NewOption(c.features[i].DisplayName, &c.features[i])
	}

	var selected *Feature
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[*Feature]().
				Title("操作を選択").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		return false
	}

	if selected.IsExit {
		return false
	}

	selected.Execute(ctx)
	fmt.Println()
	return true
}
