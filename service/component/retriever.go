package component

import (
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/sirupsen/logrus"
)

type MilvusRetrieverConfig struct {
	Collection string `yaml:"collection"`
	TopK       int    `yaml:"top_k"`
}

func NewMilvusRetriever(ctx context.Context, cfg *MilvusRetrieverConfig, cli *milvusclient.Client, emb embedding.Embedder) (*milvus2.Retriever, error) {
	retriever, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		Client:     cli,
		Collection: cfg.Collection,
		TopK:       cfg.TopK,
		SearchMode: search_mode.NewApproximate(milvus2.COSINE),
		Embedding:  emb,
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("created milvus retriever")
	return retriever, nil
}
