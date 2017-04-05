package logtee

import (
	"bytes"
	"fmt"
	"time"
)

type Event struct {
	At       time.Time
	Level    Level
	Category string
	Message  string
	Error    string
	Fields   Fields
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
				e.Fields = Fields{}
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

func (e *Event) AsKVL() string {
	if e.IsZero() {
		return ""
	}
	buff := bytes.NewBufferString("")
	fmt.Fprintf(buff, `time=%s `, escape(formatTime(e.At)))
	fmt.Fprintf(buff, `level=%s `, escape(e.Level.String()))
	fmt.Fprintf(buff, `category=%s `, escape(e.Category))
	fmt.Fprintf(buff, `msg=%s `, escape(e.Message))
	if e.Error != "" {
		fmt.Fprintf(buff, `level=%s `, escape(e.Error))
	}
	for k, v := range e.Fields {
		fmt.Fprintf(buff, `%s=%s `, k, escape(v))
	}
	return buff.String()
}

func (e *Event) String() string {
	return e.AsKVL()
}
