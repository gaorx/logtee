package logtee

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type matcher func(e *Event) (bool, error)

func matcherOf(expr string) (matcher, error) {
	if expr == "" {
		return func(e *Event) (bool, error) {
			return false, nil
		}, nil
	}
	eval, err := govaluate.NewEvaluableExpressionWithFunctions(expr, matchFuncs)
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
			"at":       e.At.Format(time.RFC3339),
			"level":    e.Level,
			"category": e.Category,
			"message":  e.Message,
			"error":    e.Error,
		}
		for k, v := range e.Fields {
			if _, ok := params[k]; !ok {
				params[k] = v
			}
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

var (
	matchFuncs = map[string]govaluate.ExpressionFunction{
		"Contains": mfContains,
	}
)

func mfContains(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, errors.New("argument count error")
	}
	s, err := strArg(args[0])
	if err != nil {
		return false, err
	}
	sub, err := strArg(args[1])
	if err != nil {
		return false, err
	}
	return strings.Contains(s, sub), nil
}
