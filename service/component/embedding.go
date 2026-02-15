package component

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/sirupsen/logrus"
)

type DashScopeEmbedderConfig struct {
	ApiKey     string `yaml:"api_key"`
	Model      string `yaml:"model"`
	Dimensions *int   `yaml:"dimensions"`
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
