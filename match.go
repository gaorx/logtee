package logtee

type Matcher func(e *Event) (bool, error)

func NewMatcher(expr string) (Matcher, error) {
	return func(e *Event) (bool, error) {
		return false, nil
	}, nil
}