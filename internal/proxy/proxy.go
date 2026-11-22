// Package proxy provides the entry point for the proxy service.
package proxy

import (
	"context"
	"errors"
)

// Options holds settings passed from the command line to the proxy service.
type Options struct {
	ConfigPath   string
	DisableColor bool
}

// Run is reserved for loading configuration and running the proxy service until cancellation.
func Run(_ context.Context, _ Options) error {
	return errors.New("proxy service is not implemented yet")
}
