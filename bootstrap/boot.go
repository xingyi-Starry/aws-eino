package bootstrap

import (
	"awseino/config"
	"awseino/lib"
	"awseino/lib/cozeloop"
	"awseino/lib/logger"
	"awseino/service/component"
	"context"

	"github.com/sirupsen/logrus"
)

func MustInit(ctx context.Context, cfg *config.Config) {
	var err error
	logger.MustInitLogger(cfg.LoggerConfig)
	lib.CozeloopClient, err = cozeloop.InitCozeloop(ctx, cfg.CozeloopConfig)
	if err != nil {
		logrus.Fatalf("Failed to init cozeloop: %v", err)
	}
	lib.Embedder, err = component.NewDashScopeEmbedder(ctx, cfg.DashScopeEmbedderConfig)
	if err != nil {
		logrus.Fatalf("Failed to create embedder: %v", err)
	}
	lib.MilvusClient, err = component.NewMilvusClient(ctx, cfg.MilvusClientConfig)
	if err != nil {
		logrus.Fatalf("Failed to create milvus client: %v", err)
	}
	lib.MilvusRetriever, err = component.NewMilvusRetriever(ctx, cfg.MilvusRetrieverConfig, lib.MilvusClient, lib.Embedder)
	if err != nil {
		logrus.Fatalf("Failed to create milvus retriever: %v", err)
	}
	lib.ChatModel, err = component.NewOpenAIChatModel(ctx, cfg.OpenAiChatModelConfig)
	if err != nil {
		logrus.Fatalf("Failed to create openai chat model: %v", err)
	}
	lib.RetrieveTool, err = component.NewRetrieveTool(ctx, lib.MilvusRetriever)
	if err != nil {
		logrus.Fatalf("Failed to create retrieve tool: %v", err)
	}
}
