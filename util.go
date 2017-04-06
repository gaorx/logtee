package logtee

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeLayoutsForParse = []string{
		time.RFC3339,
		time.RFC3339Nano,
		"20060102150405",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.Kitchen,
		"2006-01-02",                         // RFC 3339
		"2006-01-02 15:04",                   // RFC 3339 with minutes
		"2006-01-02 15:04:05",                // RFC 3339 with seconds
		"2006-01-02 15:04:05-07:00",          // RFC 3339 with seconds and timezone
		"2006-01-02T15Z0700",                 // ISO8601 with hour
		"2006-01-02T15:04Z0700",              // ISO8601 with minutes
		"2006-01-02T15:04:05Z0700",           // ISO8601 with seconds
		"2006-01-02T15:04:05.999999999Z0700", // ISO8601 with nanoseconds
	}
)

func escape(s string) string {
	return strconv.Quote(s)
}

func tryEscape(s string) string {
	if needEscape(s) {
		return escape(s)
	} else {
		return s
	}
}

func needEscape(s string) bool {
	for _, ch := range s {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func unescape(s string) (string, error) {
	return strconv.Unquote(s)
}

func parseTime(s string) (time.Time, error) {
	for _, layout := range timeLayoutsForParse {
		r, err := time.Parse(layout, s)
		if err == nil {
			return r, nil
		}
	}
	return time.Time{}, fmt.Errorf("Parse time error: %q", s)
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func lockDo(mtx *sync.Mutex, f func()) {
	defer mtx.Unlock()
	mtx.Lock()
	f()
}

func safeDo(f func()) {
	defer func() {
		if r := recover(); r != nil {
			// TODO
		}
	}()
	f()
}

func split2(s, sep string) (string, string) {
	if s == "" {
		return "", ""
	}
	ss := strings.SplitN(s, sep, 2)
	switch len(ss) {
	case 0:
		return "", ""
	case 1:
		return ss[0], ""
	default:
		return ss[0], ss[1]
	}
}

func splitNotEmpty(s, sep string) []string {
	ss := strings.Split(s, sep)
	if len(ss) == 0 {
		return ss
	}
	ss1 := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			ss1 = append(ss1, s)
		}
	}
	return ss1
}

func strArg(v interface{}) (string, error) {
	if v == nil {
		return "", errors.New("Nil arg")
	}
	if r, ok := v.(string); ok {
		return r, nil
	} else {
		return "", errors.Errorf("Argument must be string (%s)", v)
	}
}

func intArg(v interface{}) (int, error) {
	if v == nil {
		return 0, errors.New("Nil arg")
	}
	if r, ok := v.(int); ok {
		return r, nil
	} else {
		return 0, errors.Errorf("Argument must be int (%s)", v)
	}
}
