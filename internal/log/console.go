// log 包负责格式化 htunk 的控制台日志。
package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// NewConsole 创建可读的控制台日志记录器，输出 INFO 及以上级别的日志。
func NewConsole(output io.Writer, disableColor bool) *slog.Logger {
	file, isFile := output.(*os.File)
	color := !disableColor && isFile && os.Getenv("TERM") != "dumb" && term.IsTerminal(int(file.Fd()))
	return slog.New(&consoleHandler{
		output: output,
		mu:     new(sync.Mutex),
		start:  time.Now(),
		color:  color,
	})
}

type consoleHandler struct {
	output io.Writer
	mu     *sync.Mutex
	start  time.Time
	color  bool
	attrs  string
	group  string
}

func (h *consoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelInfo
}

func (h *consoleHandler) Handle(_ context.Context, record slog.Record) error {
	var line strings.Builder
	elapsed := max(0, record.Time.Sub(h.start)/time.Second)
	level := record.Level.String()
	if h.color {
		level = colorLevel(record.Level, level)
	}
	fmt.Fprintf(&line, "%s[%04d] %s", level, elapsed, record.Message)
	line.WriteString(h.attrs)
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(&line, h.group, attr)
		return true
	})
	line.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.output, line.String())
	return err
}

func (h *consoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	var suffix strings.Builder
	suffix.WriteString(h.attrs)
	for _, attr := range attrs {
		appendAttr(&suffix, h.group, attr)
	}
	next.attrs = suffix.String()
	return &next
}

func (h *consoleHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	next := *h
	next.group += name + "."
	return &next
}

func appendAttr(line *strings.Builder, group string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return
	}
	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			group += attr.Key + "."
		}
		for _, child := range attr.Value.Group() {
			appendAttr(line, group, child)
		}
		return
	}
	fmt.Fprintf(line, " %s=%q", group+attr.Key, attr.Value.String())
}

func colorLevel(level slog.Level, label string) string {
	code := "36"
	switch {
	case level >= slog.LevelError:
		code = "31"
	case level >= slog.LevelWarn:
		code = "33"
	}
	return "\x1b[" + code + "m" + label + "\x1b[0m"
}
