package logtee

import (
	"encoding/json"
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

func ParseLine(line string) (*Event, error) {
	var m map[string]string
	err := json.Unmarshal([]byte(line), &m)
	if err != nil {
		return nil, err
	}
	e := Event{}
	for n, v := range m {
		switch n {
		case "at", "time": // 'time' for logrus
			e.At, err = parseTime(v)
			if err != nil {
				return nil, err
			}
		case "level", "lvl":
			e.Level = ParseLevel(v)
		case "category":
			e.Category = v
		case "message", "msg": // 'msg' for logrus
			e.Message = v
		case "error": // 'error' for logrus
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
	b, err := jsonFormatter(e)
	if err != nil {
		return ""
	}
	return string(b)
}
