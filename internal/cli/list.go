package cli

import (
	"context"
	"fmt"
)

func (c *CLI) handleList(ctx context.Context) {
	audioList, err := c.client.List(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		return
	}

	if len(audioList) == 0 {
		fmt.Println("利用可能な音声ファイルがありません")
		return
	}

	fmt.Println("利用可能な音声:")
	for _, a := range audioList {
		fmt.Printf("  - %s\n", a)
	}
}
