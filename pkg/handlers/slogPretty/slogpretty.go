package slogpretty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdLog "log"
	"log/slog"
	"strings"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts   PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

const (
	resetCode   = "\033[0m"
	boldCode    = "\033[1m"
	magentaCode = "\033[35m"
	blueCode    = "\033[34m"
	yellowCode  = "\033[33m"
	redCode     = "\033[31m"
	cyanCode    = "\033[36m"
	whiteCode   = "\033[37m"
)


func colorize(text string, ansiCode string) string {
	if ansiCode == "" {
		return text
	}
	return fmt.Sprintf("%s%s%s", ansiCode, text, resetCode)
}

func colorToANSI(s string) string {
	
	switch {
	case s == "magenta": 
		return magentaCode
	case s == "blue":
		return blueCode
	case s == "yellow":
		return yellowCode
	case s == "red":
		return redCode
	case s == "cyan":
		return cyanCode
	case s == "white":
		return whiteCode
	default:
		return whiteCode
	}
}

var (
	magentaColor = "magenta"
	blueColor = "blue"
	yellowColor = "yellow"
	redColor = "red"
	cyanColor = "cyan"
	whiteColor = "white"
)

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = magentaColor
	case slog.LevelInfo:
		levelColor = blueColor
	case slog.LevelWarn:
		levelColor = yellowColor
	case slog.LevelError:
		levelColor = redColor
	default:
		levelColor = whiteColor
	}
	
	// Конвертируем цвет в ANSI код и применяем к логам
	level = colorize(level, colorToANSI(levelColor))

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := colorize(r.Message, colorToANSI(cyanColor))
	
	jsonStr := string(b)
	if jsonStr != "" {
		jsonStr = colorize(jsonStr, colorToANSI(whiteColor))
	}

	output := fmt.Sprintf("%s %s %s %s", 
		timeStr, 
		level, 
		msg,
		jsonStr,
	)
	
	output = strings.TrimSpace(output)
	
	h.l.Println(output)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}