//go:build go1.21

package indenthandler

import (
	"bytes"
	"regexp"
	"testing"

	"log/slog"
)

func Test(t *testing.T) {
	var buf bytes.Buffer
	l := slog.New(New(&buf, nil))
	l.Info("hello", "a", 1, "b", true, "c", 3.14, slog.Group("g", "h", 1, "i", 2), "d", "NO")
	got := buf.String()
	wantre := `time: [-0-9T:.]+Z?
level: INFO
source: ".*/indent_handler_test.go:\d+"
msg: "hello"
a: 1
b: true
c: 3.14
g:
    h: 1
    i: 2
d: "NO"
`
	re := regexp.MustCompile(wantre)
	if !re.MatchString(got) {
		t.Errorf("\ngot:\n%q\nwant:\n%q", got, wantre)
	}

	buf.Reset()
	l.Debug("test")
	if got := buf.Len(); got != 0 {
		t.Errorf("got buf.Len() = %d, want 0", got)
	}
}
