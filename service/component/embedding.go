package component

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/sirupsen/logrus"
)

type DashScopeEmbedderConfig struct {
	ApiKey     string `yaml:"api_key"`
	Model      string `yaml:"model"`
	Dimensions *int   `yaml:"dimensions"`
}

type EmbedderWrap struct {
	Embedder  embedding.Embedder
	batchSize int
}

func NewDashScopeEmbedder(ctx context.Context, config *DashScopeEmbedderConfig) (*dashscope.Embedder, error) {
	embedder, err := dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		APIKey:     config.ApiKey,
		Model:      config.Model,
		Dimensions: config.Dimensions,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("created embedder: %s", config.Model)
	return embedder, nil
}

func NewEmbedderWrap(ctx context.Context, embedder embedding.Embedder, batchSize int) *EmbedderWrap {
	return &EmbedderWrap{
		Embedder:  embedder,
		batchSize: batchSize,
	}
}

func (e *EmbedderWrap) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	embedStrings := make([][]float64, 0, len(texts))
	for i := 0; i < len(texts); i += e.batchSize {
		batch := texts[i:min(i+e.batchSize, len(texts))]
		res, err := e.Embedder.EmbedStrings(ctx, batch, opts...)
		if err != nil {
			return nil, err
		}
		embedStrings = append(embedStrings, res...)
	}

	return embedStrings, nil
}
