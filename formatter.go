package logtee

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
)

type Formatter func(*Event) ([]byte, error)
type FormatterFactory func(Config) (Formatter, error)

var (
	formatterFactories = map[string]FormatterFactory{}
)

func RegisterFormatterFactory(name string, ff FormatterFactory) {
	if name == "" {
		return
	}
	if ff != nil {
		formatterFactories[name] = ff
	} else {
		delete(formatterFactories, name)
	}
}

func CompileFormatter(conf Config) (Formatter, error) {
	formatType := conf.Str("name", "")
	if formatType == "" {
		return nil, errors.New("Nil format")
	}
	ff, _ := formatterFactories[formatType]
	if ff == nil {
		return nil, errors.Errorf("Not found format: %s", formatType)
	}
	return ff(conf)
}

func FormatterFactoryOf(f Formatter) FormatterFactory {
	return func(_ Config) (Formatter, error) {
		return f, nil
	}
}

func jsonFormatter(e *Event) ([]byte, error) {
	m := map[string]string{
		"time":     formatTime(e.At),
		"level":    e.Level.String(),
		"category": e.Category,
		"msg":      e.Message,
	}
	if e.Error != "" {
		m["error"] = e.Error
	}
	for k, v := range e.Fields {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

func newCsvFormatter(conf Config) (Formatter, error) {
	fields := splitNotEmpty(conf.Str("fields", ""), ",")
	nFields := len(fields)
	sep := conf.Str("sep", "\t")
	return func(e *Event) ([]byte, error) {
		buff := bytes.NewBufferString("")
		for i, field := range fields {
			switch field {
			case "at":
				buff.WriteString(tryEscape(formatTime(e.At)))
			case "level":
				buff.WriteString(e.Level.String())
			case "category":
				buff.WriteString(tryEscape(e.Category))
			case "message":
				buff.WriteString(tryEscape(e.Message))
			case "error":
				buff.WriteString(tryEscape(e.Error))
			default:
				v, _ := e.Fields[field]
				buff.WriteString(tryEscape(v))
			}
			if i < nFields-1 {
				buff.WriteString(sep)
			}
		}
		return buff.Bytes(), nil
	}, nil
}
