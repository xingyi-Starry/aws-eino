package compose

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type ChainConfig struct {
	ChatModel model.ToolCallingChatModel
	Tools     []tool.BaseTool
}

func NewChain(ctx context.Context, cfg *ChainConfig) (compose.Runnable[string, string], error) {
	chain := compose.NewChain[string, string]()

	lbd1 := compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
		return map[string]any{
			"input": input,
		}, nil
	})

	lbd2 := compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
		builder := strings.Builder{}
		for _, msg := range input {
			builder.WriteString(msg.Content)
			builder.WriteString("\n")
		}
		return builder.String(), nil
	})

	tpl := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一名知识库问答助手，请根据用户的问题检索知识库。"),
		schema.UserMessage("{input}"),
	)

	toolInfos := make([]*schema.ToolInfo, 0, len(cfg.Tools))
	for _, tool := range cfg.Tools {
		toolInfo, err := tool.Info(ctx)
		if err != nil {
			logrus.Errorf("failed to get tool info: %v", err)
			return nil, err
		}
		toolInfos = append(toolInfos, toolInfo)
	}
	cm, err := cfg.ChatModel.WithTools(toolInfos)
	if err != nil {
		logrus.Errorf("failed to create chat model: %v", err)
		return nil, err
	}

	toolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: cfg.Tools,
	})
	if err != nil {
		logrus.Errorf("failed to create tool node: %v", err)
		return nil, err
	}

	chain.AppendLambda(lbd1)
	chain.AppendChatTemplate(tpl)
	chain.AppendChatModel(cm)
	chain.AppendToolsNode(toolNode)
	chain.AppendLambda(lbd2)

	return chain.Compile(ctx)
}
