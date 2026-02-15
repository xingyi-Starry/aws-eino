package component

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/sirupsen/logrus"
)

type OpenAiChatModelConfig struct {
	BaseUrl string `yaml:"base_url"`
	ApiKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

func NewOpenAIChatModel(ctx context.Context, config *OpenAiChatModelConfig) (*openai.ChatModel, error) {
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: config.BaseUrl,
		APIKey:  config.ApiKey,
		Model:   config.Model,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("created chat model: %s", config.Model)
	return cm, nil
}
