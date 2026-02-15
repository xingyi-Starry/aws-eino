package main

import (
	"awseino/config"
	"awseino/service/component"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	embedder, err := component.NewDashScopeEmbedder(ctx, cfg.DashScopeEmbedderConfig)
	if err != nil {
		logrus.Fatalf("Failed to create embedder: %v", err)
	}

	res, err := embedder.EmbedStrings(ctx, []string{"你说得对，但是机器学习是人工智能的一个分支，它是一种让计算机系统能够从数据中学习的技术。", "陈墨站在公寓窗前，看着手机屏幕上的倒计时归零。窗外，这座拥有八百万人口的城市正在沉睡——或者说，假装沉睡。街道上的路灯在凌晨的薄雾中晕开昏黄的光圈，几辆清洁车缓慢驶过，发出机械的嗡鸣。一切如常，除了他胸腔里那颗跳得过于用力的心脏。"})
	if err != nil {
		logrus.Fatalf("Failed to embed strings: %v", err)
	}

	fmt.Println("length of res:", len(res))
	fmt.Println("vector dim:", len(res[0]))
}
