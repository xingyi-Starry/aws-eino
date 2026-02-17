package compose

import (
	"context"
	"errors"
	"io"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type ReActConfig struct {
	ChatModel    model.ToolCallingChatModel
	Tools        []tool.BaseTool
	SystemPrompt string
}

type reActState struct {
	histories []*schema.Message
}

func genReActState(ctx context.Context) *reActState { return &reActState{} }

func addHistory(ctx context.Context, in *schema.Message, state *reActState) (*schema.Message, error) {
	state.histories = append(state.histories, in)
	return in, nil
}

func addHistories(ctx context.Context, in []*schema.Message, state *reActState) ([]*schema.Message, error) {
	state.histories = append(state.histories, in...)
	return state.histories, nil
}

const (
	nodeTpl       = "tpl"
	nodeChatModel = "chat_model"
	nodeBranch    = "branch"
	nodeToolCall  = "tool_call"
	nodeOutput    = "output"
)

func NewReAct(ctx context.Context, cfg *ReActConfig) (compose.Runnable[map[string]any, string], error) {
	g := compose.NewGraph[map[string]any, string](compose.WithGenLocalState(genReActState))

	tpl := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(cfg.SystemPrompt),
		schema.MessagesPlaceholder("history", true),
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
		logrus.Errorf("failed to bind tools to chat model: %v", err)
		return nil, err
	}

	branch := compose.NewStreamGraphBranch(func(ctx context.Context, input *schema.StreamReader[*schema.Message]) (endNode string, err error) {
		count := 0
		defer input.Close()
		for {
			msg, err := input.Recv()
			count++
			if err == io.EOF {
				return nodeOutput, nil
			}
			if err != nil {
				return nodeOutput, err
			}

			if len(msg.ToolCalls) > 0 {
				return nodeToolCall, nil
			}

			if len(msg.Content) == 0 { // skip empty chunks at the front
				continue
			}

			return nodeOutput, nil
		}
	}, map[string]bool{
		nodeToolCall: true,
		nodeOutput:   true,
	})

	toolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{Tools: cfg.Tools})
	if err != nil {
		logrus.Errorf("failed to create tool node: %v", err)
		return nil, err
	}

	output := compose.TransformableLambda(
		func(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.StreamReader[string], error) {
			return schema.StreamReaderWithConvert(input, func(src *schema.Message) (string, error) {
				return src.Content, nil
			}), nil
		},
	)

	eh := &errHandler{}
	eh.handleErr(g.AddChatTemplateNode(nodeTpl, tpl))
	eh.handleErr(g.AddChatModelNode(nodeChatModel, cm, compose.WithStatePreHandler(addHistories)))
	eh.handleErr(g.AddToolsNode(nodeToolCall, toolNode, compose.WithStatePreHandler(addHistory)))
	eh.handleErr(g.AddLambdaNode(nodeOutput, output))
	eh.handleErr(g.AddBranch(nodeChatModel, branch))
	if err = eh.GetError(); err != nil {
		logrus.Errorf("failed to add nodes to graph: %v", err)
		return nil, err
	}
	eh.Clear()

	eh.handleErr(g.AddEdge(compose.START, nodeTpl))
	eh.handleErr(g.AddEdge(nodeTpl, nodeChatModel))
	eh.handleErr(g.AddEdge(nodeToolCall, nodeChatModel))
	eh.handleErr(g.AddEdge(nodeOutput, compose.END))
	if err = eh.GetError(); err != nil {
		logrus.Errorf("failed to add edges to graph: %v", err)
		return nil, err
	}

	return g.Compile(ctx)
}

type errHandler struct {
	errs []error
}

func (e *errHandler) handleErr(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

func (e *errHandler) GetError() error {
	return errors.Join(e.errs...)
}

func (e *errHandler) Error() string {
	err := e.GetError()
	return err.Error()
}

func (e *errHandler) Clear() {
	e.errs = nil
}
