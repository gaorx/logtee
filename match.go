package logtee

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	_ "github.com/Knetic/govaluate"
)

type matcher func(e *Event) (bool, error)

func matcherOf(expr string) (matcher, error) {
	if expr == "" {
		return func(e *Event) (bool, error) {
			return false, nil
		}, nil
	}
	eval, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return nil, err
	}
	return func(e *Event) (bool, error) {
		if e == nil {
			return false, errors.New("Nil event")
		}
		params := map[string]interface{}{
			// levels
			"BIZ":   BizLevel,
			"PANIC": PanicLevel,
			"FATAL": FatalLevel,
			"ERROR": ErrorLevel,
			"WARN":  WarnLevel,
			"INFO":  InfoLevel,
			"DEBUG": DebugLevel,

			// major
			"time":     e.At,
			"level":    e.Level,
			"category": e.Category,
			"msg":      e.Message,
			"err":      e.Error,
			"fields":   e.Fields,
		}
		res, err := eval.Evaluate(params)
		if err != nil {
			return false, err
		}
		b, ok := res.(bool)
		if !ok {
			return false, fmt.Errorf("The type of result is not bool (%v)", res)
		}
		return b, nil
	}, nil
}
