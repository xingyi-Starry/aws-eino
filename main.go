package main

import (
	"awseino/config"
	"awseino/service/component"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino/schema"
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

	var input string
	fmt.Print("> ")
	fmt.Scanln(&input)

	docs, err := retriever.Retrieve(ctx, input)
	if err != nil {
		logrus.Fatalf("Failed to retrieve documents: %v", err)
	}

	var sysMsgContent strings.Builder
	sysMsgContent.WriteString("以下是一些相关文档，请根据文档内容回答问题。")
	for _, doc := range docs {
		fmt.Fprintf(&sysMsgContent, "\n\n---\n\n%s", doc.Content)
		logrus.Infof("retrieved document: %s", doc.ID)
	}

	msgs := []*schema.Message{
		schema.SystemMessage(sysMsgContent.String()),
		schema.UserMessage(input),
	}

	stream, err := cm.Stream(ctx, msgs)
	fmt.Println("")
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				logrus.Errorf("Failed to receive chunk: %v", err)
			}
			break
		}
		fmt.Print(chunk.Content)
	}
	logrus.Infof("stream finished")
}
