//go:build go1.21 && !race

package indenthandler

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"testing"
)

func TestAlloc(t *testing.T) {
	a := slog.String("key", "value")
	t.Run("Appendf", func(t *testing.T) {
		buf := make([]byte, 0, 100)
		g := testing.AllocsPerRun(2, func() {
			buf = fmt.Appendf(buf, "%s: %q\n", a.Key, a.Value.String())
		})
		if g, w := int(g), 2; g != w {
			t.Errorf("got %d, want %d", g, w)
		}
	})
	t.Run("appends", func(t *testing.T) {
		buf := make([]byte, 0, 100)
		g := testing.AllocsPerRun(2, func() {
			buf = append(buf, a.Key...)
			buf = append(buf, ": "...)
			buf = strconv.AppendQuote(buf, a.Value.String())
			buf = append(buf, '\n')
		})
		if g, w := int(g), 0; g != w {
			t.Errorf("got %d, want %d", g, w)
		}
	})

	t.Run("Handle", func(t *testing.T) {
		l := slog.New(New(io.Discard, nil))
		got := testing.AllocsPerRun(10, func() {
			l.LogAttrs(nil, slog.LevelInfo, "hello", slog.String("a", "1"))
		})
		if g, w := int(got), 6; g > w {
			t.Errorf("got %d, want at most %d", g, w)
		}
	})
}
