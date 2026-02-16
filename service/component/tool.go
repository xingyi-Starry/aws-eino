package component

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type RetrieveTool struct {
	retriever retriever.Retriever
}

type RetrieveToolParams struct {
	Query string `json:"query" jsonschema_description:"content to search"`
}

type RetrieveToolResult struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (t *RetrieveTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	paramSchema, err := toolutils.GoStruct2ParamsOneOf[RetrieveToolParams]()
	if err != nil {
		logrus.Errorf("failed to convert retrieve tool params to json schema: %+v", err)
		return nil, err
	}
	return &schema.ToolInfo{
		Name:        "retriever",
		Desc:        "Retrieve documents from the knowledge base",
		ParamsOneOf: paramSchema,
	}, nil
}

func (t *RetrieveTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params RetrieveToolParams
	if err := sonic.Unmarshal([]byte(arguments), &params); err != nil {
		return "", fmt.Errorf("failed to unmarshal arguments: %+v", err)
	}
	logrus.Infof("retrieving documents for query: [%s]", params.Query)
	docs, err := t.retriever.Retrieve(ctx, params.Query)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve documents: %+v", err)
	}

	results := make([]RetrieveToolResult, 0, len(docs))
	for _, doc := range docs {
		results = append(results, RetrieveToolResult{
			ID:      doc.ID,
			Content: doc.Content,
		})
	}
	return sonic.MarshalString(results)
}

func NewRetrieveTool(ctx context.Context, retriever retriever.Retriever) (tool.InvokableTool, error) {
	if retriever == nil {
		return nil, fmt.Errorf("retriever is nil")
	}
	tool := &RetrieveTool{
		retriever: retriever,
	}
	return tool, nil
}
