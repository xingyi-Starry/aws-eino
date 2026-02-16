package common

import (
	"context"

	"github.com/cloudwego/eino/callbacks"
	"github.com/sirupsen/logrus"
)

func GenCallback() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		logrus.Debugf("当前%s节点输入:%s\n", info.Component, input)
		return ctx
	}).OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		logrus.Debugf("当前%s节点输出:%s\n", info.Component, output)
		return ctx
	}).Build()
	return handler
}
