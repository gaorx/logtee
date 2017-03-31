package logtee

import (
	"fmt"
	"strconv"
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
