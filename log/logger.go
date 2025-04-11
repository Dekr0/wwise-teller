package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
	"time"
)

type Handler struct {
	opts Options
	/* TODO: WithAttr and WithGroup */
	/*
	 Ensure atomic write to io.Writer.
	 Use of pointer make sure WithGroup and WithAttrs use the same mutex since 
	 they copy the handler.
	*/
	mu *sync.Mutex 
	out io.Writer
}

type Options struct {
	Level slog.Leveler
}

func NewHandler(out io.Writer, opts *Options) *Handler {
	h := &Handler{ out: out, mu: &sync.Mutex{} }
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

/*
 Handler receive a `slog.Record`, process it, and either write it to an io.Write,
 or send it to the next handler.
*/
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	/* 
	 Allocate large enough buffer for most single record upfront.
	 Avoid copying and allocation that happen when the inital slice is empty or 
	 small.
	*/
	buf := make([]byte, 0, 1024) 

	/* Ignore zero time (Common rules) */
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}

	/* Level */
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)

	/* AddSource */
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r. PC})
		f, _ := fs.Next()
		buf = h.appendAttr(
			buf,
			slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)),
			0,
		)
	}

	/* Message */
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)

	/*
	indentLevel := 0
	r.Attrs(func (a slog.Attr) bool {
		buf = h.appendAttr(buf, a, indentLevel)
		return true
	})
	*/

	buf = append(buf, "---\n"...)

	/* Atomic write */
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)

	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}

func (h *Handler) appendAttr(buf []byte, a slog.Attr, indentLevel int) []byte {
	a.Value = a.Value.Resolve()

	/* Ignore empty attributes (Common rules) */
	if a.Equal(slog.Attr{}) {
		return buf
	}

	/* Indent 4 space per level */
	buf = fmt.Appendf(buf, "%*s", indentLevel * 4, "")

	switch a.Value.Kind() {
	case slog.KindString:
		/* Quoted string */
		buf = fmt.Appendf(buf, "%s: %q\n", a.Key, a.Value.String())
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s: %s\n", a.Key, a.Value.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		attrs := a.Value.Group()
		/* Group with no attributes is ignored (Common rules) */
		if len(attrs) == 0 {
			return buf
		}
		/*
		 If the key is non-empty, write it out and indent the rest of attrs.
		 Otherwise, inline the attrs
		*/
		if a.Key != "" {
			buf = fmt.Appendf(buf, "%s:\n", a.Key)
			indentLevel++
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga, indentLevel)
		}
	default:
		buf = fmt.Appendf(buf, "%s: %s\n", a.Key, a.Value)
	}
	return buf
}
