package lib

import (
	mcIndexer "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	mcRetriever "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var Embedder embedding.Embedder
var MilvusClient *milvusclient.Client
var MilvusIndexer *mcIndexer.Indexer
var MilvusRetriever *mcRetriever.Retriever
var ChatModel model.ToolCallingChatModel

var RetrieveTool tool.InvokableTool
