package cozeloop

import (
	"context"

	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
)

type CozeloopConfig struct {
	ApiToken    string `yaml:"api_token"`
	WorkspaceID string `yaml:"workspace_id"`
}

func InitCozeloop(ctx context.Context, cfg *CozeloopConfig) (cozeloop.Client, error) {
	client, err := cozeloop.NewClient(
		cozeloop.WithWorkspaceID(cfg.WorkspaceID),
		cozeloop.WithAPIToken(cfg.ApiToken),
	)
	if err != nil {
		return nil, err
	}
	cozeloop.SetDefaultClient(client)
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	return client, nil
}
