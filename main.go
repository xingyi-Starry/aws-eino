package main

import (
	"awseino/bootstrap"
	"awseino/config"
	"awseino/lib"
	"awseino/service/common"
	cps "awseino/service/compose"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

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

	defer lib.CozeloopClient.Close(ctx)

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

	output, err := react.Stream(ctx, map[string]any{"input": input}, compose.WithCallbacks(common.GenCallback()))
	if err != nil {
		logrus.Fatalf("Failed to invoke react: %v", err)
	}

	for {
		chunk, err := output.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				logrus.Errorf("Failed to receive chunk: %v", err)
			}
			logrus.Infof("Stream ended")
			break
		}
		fmt.Print(chunk)
	}

	time.Sleep(5 * time.Second) // 太nm荒谬了，cozeloop没有一个可用的异步等待上报完成的接口
}
