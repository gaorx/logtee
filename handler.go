package logtee

import (
	"errors"
	"fmt"
	"strings"
)

type Handler interface {
	Init() error
	Close() error
	Match(e *Event) (bool, error)
	Append(e *Event)
	Set(k string, v interface{})
	Get(k string, def interface{}) interface{}
}

type Handlers []Handler

type HandlerFactory func(name string, conf Config) Handler

var (
	handlerFactories map[string]HandlerFactory = map[string]HandlerFactory{}
)

func RegisterHandlerFactory(typ string, factory HandlerFactory) {
	if factory != nil {
		handlerFactories[typ] = factory
	} else {
		delete(handlerFactories, typ)
	}
}

func NewHandler(typ, name string, conf Config) (Handler, error) {
	factory, ok := handlerFactories[typ]
	if !ok || factory == nil {
		return nil, fmt.Errorf("Not found handler type:%s", typ)
	}
	h := factory(name, conf)
	if h == nil {
		return nil, errors.New("Nil handler")
	}
	return h, nil
}

func ParseHandlers(conf Config) (Handlers, error) {
	var handlers Handlers
	for k := range conf {
		typ, name := split2(k, ":")
		h, err := NewHandler(typ, name, conf.Sub(k, nil))
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, h)
	}
	return handlers, nil
}

func (hh Handlers) Init() error {
	var inited Handlers
	for _, h := range hh {
		err := h.Init()
		if err != nil {
			inited.Close()
			return err
		}
		inited = append(inited, h)
	}
	return nil
}

func (hh Handlers) Close() error {
	for _, h := range hh {
		h.Close()
	}
	return nil
}

func (hh Handlers) DoOne(e *Event) {
	var matched []Handler
	for _, h := range hh {
		b, err := h.Match(e)
		if err != nil {
			// TODO: log
			b = false
		}
		if b {
			matched = append(matched, h)
		}
	}
	for _, h := range matched {
		h.Append(e)
	}
}

func (hh Handlers) Do(src Source) {
	if src == nil {
		return
	}
	for line := range src {
		line = strings.TrimRight(line, "\r\n ")
		e, err := ParseKVL(line)
		if err != nil {
			continue
		}
		hh.DoOne(e)
	}
}
