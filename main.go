package main

import (
	"awseino/bootstrap"
	"awseino/config"
	"awseino/lib"
	"awseino/service/common"
	cps "awseino/service/compose"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	bootstrap.MustInit(ctx, cfg)

	react, err := cps.NewReAct(ctx, &cps.ReActConfig{
		ChatModel:    lib.ChatModel,
		Tools:        []tool.BaseTool{lib.RetrieveTool},
		SystemPrompt: "你是一名知识库问答助手，请根据用户的问题检索知识库并回答用户的问题。",
	})
	if err != nil {
		logrus.Fatalf("Failed to create react: %v", err)
	}

	var input string
	fmt.Print("> ")
	fmt.Scanln(&input)

	output, err := react.Invoke(ctx, map[string]any{"input": input}, compose.WithCallbacks(common.GenCallback()))
	if err != nil {
		logrus.Fatalf("Failed to invoke react: %v", err)
	}
	fmt.Println(output)
}
