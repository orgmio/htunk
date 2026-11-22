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

	"github.com/orgmio/htunk/internal/config"
	"golang.org/x/term"
)

// NewConsole 创建用于报告配置加载错误的默认日志记录器。
func NewConsole(output io.Writer, disableColor bool) *slog.Logger {
	return newLogger(output, slog.LevelInfo, disableColor, false, false)
}

// Open 根据配置创建日志记录器，并返回释放输出资源的函数。
func Open(options config.Log, stderr io.Writer, disableColor bool) (*slog.Logger, func() error, error) {
	closeOutput := func() error { return nil }
	if !options.Enable {
		return slog.New(slog.DiscardHandler), closeOutput, nil
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(options.Level)); err != nil {
		return nil, nil, err
	}

	output := stderr
	toFile := options.Output != ""
	if toFile {
		file, err := os.OpenFile(options.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			return nil, nil, fmt.Errorf("open log file: %w", err)
		}
		output = file
		closeOutput = file.Close
	}
	return newLogger(output, level, disableColor, options.Timestamp, toFile), closeOutput, nil
}

func newLogger(output io.Writer, level slog.Level, disableColor, fullTimestamp, toFile bool) *slog.Logger {
	file, isFile := output.(*os.File)
	color := !disableColor && !toFile && isFile && os.Getenv("TERM") != "dumb" && term.IsTerminal(int(file.Fd()))
	return slog.New(&consoleHandler{
		output:        output,
		mu:            new(sync.Mutex),
		start:         time.Now(),
		level:         level,
		color:         color,
		fullTimestamp: fullTimestamp,
		elapsed:       !toFile,
	})
}

type consoleHandler struct {
	output        io.Writer
	mu            *sync.Mutex
	start         time.Time
	level         slog.Level
	color         bool
	fullTimestamp bool
	elapsed       bool
	attrs         string
	group         string
}

func (h *consoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *consoleHandler) Handle(_ context.Context, record slog.Record) error {
	var line strings.Builder
	level := record.Level.String()
	if h.color {
		level = colorLevel(record.Level, level)
	}
	switch {
	case h.fullTimestamp:
		fmt.Fprintf(&line, "%s %s ", record.Time.Format("-0700 2006-01-02 15:04:05"), level)
	case h.elapsed:
		elapsed := max(0, record.Time.Sub(h.start)/time.Second)
		fmt.Fprintf(&line, "%s[%04d] ", level, elapsed)
	default:
		line.WriteString(level)
		line.WriteByte(' ')
	}
	line.WriteString(strings.TrimSuffix(record.Message, "\n"))
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
	case level < slog.LevelInfo:
		code = "37"
	}
	return "\x1b[" + code + "m" + label + "\x1b[0m"
}
