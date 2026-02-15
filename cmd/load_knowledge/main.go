package main

import (
	"awseino/config"
	"awseino/service/component"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	embedderRaw, err := component.NewDashScopeEmbedder(ctx, cfg.DashScopeEmbedderConfig)
	if err != nil {
		logrus.Fatalf("failed to create embedder: %v", err)
	}
	embedder := component.NewEmbedderWrap(ctx, embedderRaw, 25)

	mc, err := component.NewMilvusClient(ctx, cfg.MilvusClientConfig)
	if err != nil {
		logrus.Fatalf("failed to create milvus client: %v", err)
	}

	indexer, err := component.NewMilvusIndexer(ctx, cfg.MilvusIndexerConfig, mc, embedder)
	if err != nil {
		logrus.Fatalf("failed to create indexer: %v", err)
	}

	loader, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true,
	})
	if err != nil {
		logrus.Fatalf("failed to create file loader: %v", err)
	}

	trans, err := semantic.NewSplitter(ctx, &semantic.Config{
		Embedding:    embedder,
		BufferSize:   4,
		MinChunkSize: 300,
		Separators:   []string{"\r\n\r\n", "\n\n", "***", "---", "==="},
		Percentile:   0.8,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return fmt.Sprintf("%s-chunk-%d", originalID, splitIndex)
		},
	})
	if err != nil {
		logrus.Fatalf("failed to create semantic splitter: %v", err)
	}

	files, err := os.ReadDir("./data")
	if err != nil {
		logrus.Fatalf("failed to read data directory: %v", err)
	}

	docs := make([]*schema.Document, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			doc, err := loader.Load(ctx, document.Source{
				URI: path.Join("./data", file.Name()),
			})
			if err != nil {
				logrus.Errorf("failed to load file %s: %v", file.Name(), err)
				continue
			}

			docs = append(docs, doc...)
		}
	}

	logrus.Infof("transforming %d documents", len(docs))
	transDocs, err := trans.Transform(ctx, docs)
	if err != nil {
		logrus.Fatalf("failed to transform documents: %v", err)
	}
	logrus.Infof("transformed %d documents", len(transDocs))

	logrus.Infof("storing %d documents", len(transDocs))
	ids, err := indexer.Store(ctx, transDocs)
	if err != nil {
		logrus.Fatalf("failed to store documents: %v", err)
	}
	fmt.Println("stored documents with ids:")
	for _, id := range ids {
		fmt.Println(id)
	}
}
