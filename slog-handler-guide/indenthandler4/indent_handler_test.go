//go:build go1.21

package indenthandler

import (
	"bytes"
	"log/slog"
	"reflect"
	"regexp"
	"testing"
	"testing/slogtest"

	"gopkg.in/yaml.v3"
)

// !+TestSlogtest
func TestSlogtest(t *testing.T) {
	var buf bytes.Buffer
	err := slogtest.TestHandler(New(&buf, nil), func() []map[string]any {
		return parseLogEntries(t, buf.Bytes())
	})
	if err != nil {
		t.Error(err)
	}
}

// !-TestSlogtest

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

// !+parseLogEntries
func parseLogEntries(t *testing.T, data []byte) []map[string]any {
	entries := bytes.Split(data, []byte("---\n"))
	entries = entries[:len(entries)-1] // last one is empty
	var ms []map[string]any
	for _, e := range entries {
		var m map[string]any
		if err := yaml.Unmarshal([]byte(e), &m); err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}
	return ms
}

// !-parseLogEntries

func TestParseLogEntries(t *testing.T) {
	in := `
a: 1
b: 2
c: 3
g:
    h: 4
    i: five
d: 6
---
e: 7
---
`
	want := []map[string]any{
		{
			"a": 1,
			"b": 2,
			"c": 3,
			"g": map[string]any{
				"h": 4,
				"i": "five",
			},
			"d": 6,
		},
		{
			"e": 7,
		},
	}
	got := parseLogEntries(t, []byte(in[1:]))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n%v\nwant:\n%v", got, want)
	}
}
