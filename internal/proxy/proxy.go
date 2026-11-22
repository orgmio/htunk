// proxy 包提供代理服务的运行入口。
package proxy

import (
	"context"
	"errors"
)

// Options 保存命令行传给代理服务的设置。
type Options struct {
	ConfigPath string
}

// Run 预留用于加载配置并运行代理服务，直到上下文取消。
func Run(_ context.Context, _ Options) error {
	return errors.New("proxy service is not implemented yet")
}
