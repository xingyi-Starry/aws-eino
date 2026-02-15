package component

import (
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/sirupsen/logrus"
)

type MilvusIndexerConfig struct {
	Address    string `yaml:"address"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Collection string `yaml:"collection"`
	Dimension  int64  `yaml:"dimension"`
}

func NewMilvusClient(ctx context.Context, config *MilvusIndexerConfig) (*milvusclient.Client, error) {
	//初始化客户端
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: config.Address,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("created milvus client")
	return client, nil
}

func NewMilvusIndexer(ctx context.Context, config *MilvusIndexerConfig, embedder embedding.Embedder) (*milvus2.Indexer, error) {
	client, err := NewMilvusClient(ctx, config)
	if err != nil {
		return nil, err
	}

	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		Client:     client,
		Collection: config.Collection,
		Vector: &milvus2.VectorConfig{
			Dimension:    config.Dimension,
			MetricType:   milvus2.COSINE,
			IndexBuilder: milvus2.NewHNSWIndexBuilder().WithM(16).WithEfConstruction(200),
		},
		Embedding: embedder,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("created milvus indexer")
	return indexer, nil
}
