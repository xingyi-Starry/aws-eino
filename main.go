package main

import (
	"awseino/config"
	"awseino/service/component"
	"awseino/service/compose"
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/tool"
	"github.com/sirupsen/logrus"
)

func main() {
	logFile, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(logFile)

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	embedder, err := component.NewDashScopeEmbedder(ctx, cfg.DashScopeEmbedderConfig)
	if err != nil {
		logrus.Fatalf("Failed to create embedder: %v", err)
	}
	mc, err := component.NewMilvusClient(ctx, cfg.MilvusClientConfig)
	if err != nil {
		logrus.Fatalf("Failed to create milvus client: %v", err)
	}
	retriever, err := component.NewMilvusRetriever(ctx, cfg.MilvusRetrieverConfig, mc, embedder)
	if err != nil {
		logrus.Fatalf("Failed to create milvus retriever: %v", err)
	}
	cm, err := component.NewOpenAIChatModel(ctx, cfg.OpenAiChatModelConfig)
	if err != nil {
		logrus.Fatalf("Failed to create openai chat model: %v", err)
	}
	retrieveTool, err := component.NewRetrieveTool(ctx, retriever)
	if err != nil {
		logrus.Fatalf("Failed to create retrieve tool: %v", err)
	}

	chain, err := compose.NewChain(ctx, &compose.ChainConfig{
		ChatModel: cm,
		Tools:     []tool.BaseTool{retrieveTool},
	})
	if err != nil {
		logrus.Fatalf("Failed to create chain: %v", err)
	}

	var input string
	fmt.Print("> ")
	fmt.Scanln(&input)

	output, err := chain.Invoke(ctx, input)
	if err != nil {
		logrus.Fatalf("Failed to invoke chain: %v", err)
	}
	fmt.Println(output)
}
