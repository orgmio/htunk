// proxy 包提供代理服务的运行入口。
package proxy

import (
	"context"
	"errors"

	"github.com/orgmio/htunk/internal/config"
)

// Run 预留用于使用已加载的配置运行代理服务，直到上下文取消。
func Run(_ context.Context, _ config.Config) error {
	return errors.New("proxy service is not implemented yet")
}
