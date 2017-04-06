package logtee

import (
	"time"
)

type Event struct {
	At       time.Time
	Level    Level
	Category string
	Message  string
	Error    string
	Fields   map[string]string
}

func ParseKVL(line string) (*Event, error) {
	fieldEntries, err := parseKvlEntries(line)
	if err != nil {
		return nil, err
	}
	e := Event{}
	for _, fieldEntry := range fieldEntries {
		n, v := fieldEntry.n, fieldEntry.v
		switch n {
		case "time":
			e.At, err = parseTime(v)
			if err != nil {
				return nil, err
			}
		case "level", "lvl":
			e.Level = ParseLevel(v)
		case "category":
			e.Category = v
		case "message", "msg":
			e.Message = v
		case "error", "err":
			e.Error = v
		default:
			if e.Fields == nil {
				e.Fields = map[string]string{}
			}
			e.Fields[n] = v
		}
	}
	return &e, nil
}

func (e *Event) IsZero() bool {
	if e == nil {
		return true
	}
	return e.At.IsZero() &&
		e.Level == UnknownLevel &&
		e.Category == "" &&
		e.Message == "" &&
		e.Error == "" &&
		e.Fields == nil
}

func (e *Event) String() string {
	b, err := kvlFormatter(e)
	if err != nil {
		return ""
	}
	return string(b)
}
