package config

import (
	"awseino/lib/cozeloop"
	"awseino/lib/logger"
	"awseino/service/component"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LoggerConfig            *logger.LogConfig                  `yaml:"logger"`
	OpenAiChatModelConfig   *component.OpenAiChatModelConfig   `yaml:"open_ai_chat_model"`
	DashScopeEmbedderConfig *component.DashScopeEmbedderConfig `yaml:"dash_scope_embedder"`
	MilvusClientConfig      *component.MilvusClientConfig      `yaml:"milvus_client"`
	MilvusIndexerConfig     *component.MilvusIndexerConfig     `yaml:"milvus_indexer"`
	MilvusRetrieverConfig   *component.MilvusRetrieverConfig   `yaml:"milvus_retriever"`
	CozeloopConfig          *cozeloop.CozeloopConfig           `yaml:"cozeloop"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
