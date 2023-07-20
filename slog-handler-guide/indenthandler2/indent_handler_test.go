//go:build go1.21

package indenthandler

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"log/slog"
	"testing/slogtest"
)

func TestSlogtest(t *testing.T) {
	var buf bytes.Buffer
	err := slogtest.TestHandler(New(&buf, nil), func() []map[string]any {
		return parseLogEntries(buf.String())
	})
	if err != nil {
		t.Error(err)
	}
}

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

func parseLogEntries(s string) []map[string]any {
	var ms []map[string]any
	scan := bufio.NewScanner(strings.NewReader(s))
	for scan.Scan() {
		m := parseGroup(scan)
		ms = append(ms, m)
	}
	if scan.Err() != nil {
		panic(scan.Err())
	}
	return ms
}

func parseGroup(scan *bufio.Scanner) map[string]any {
	m := map[string]any{}
	groupIndent := -1
	for {
		line := scan.Text()
		if line == "---" { // end of entry
			break
		}
		k, v, found := strings.Cut(line, ":")
		if !found {
			panic(fmt.Sprintf("no ':' in line %q", line))
		}
		indent := strings.IndexFunc(k, func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		if indent < 0 {
			panic("blank line")
		}
		if groupIndent < 0 {
			// First line in group; remember the indent.
			groupIndent = indent
		} else if indent < groupIndent {
			// End of group
			break
		} else if indent > groupIndent {
			panic(fmt.Sprintf("indent increased on line %q", line))
		}

		key := strings.TrimSpace(k)
		if v == "" {
			// Just a key: start of a group.
			if !scan.Scan() {
				panic("empty group")
			}
			m[key] = parseGroup(scan)
		} else {
			v = strings.TrimSpace(v)
			if len(v) > 0 && v[0] == '"' {
				var err error
				v, err = strconv.Unquote(v)
				if err != nil {
					panic(err)
				}
			}
			m[key] = v
			if !scan.Scan() {
				break
			}
		}
	}
	return m
}

func TestParseLogEntries(t *testing.T) {
	in := `
a: 1
b: 2
c: 3
g:
    h: 4
    i: 5
d: 6
---
e: 7
---
`
	want := []map[string]any{
		{
			"a": "1",
			"b": "2",
			"c": "3",
			"g": map[string]any{
				"h": "4",
				"i": "5",
			},
			"d": "6",
		},
		{
			"e": "7",
		},
	}
	got := parseLogEntries(in[1:])
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:\n%v\nwant:\n%v", got, want)
	}
}
