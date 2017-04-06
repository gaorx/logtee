package logtee

import (
	"github.com/pkg/errors"
	"sync"
	"time"
)

type BaseHandler struct {
	Name          string
	Config        Config
	Processor     func(events []*Event) error
	Formatter     Formatter
	matcher       matcher
	buff          []*Event
	buffMtx       sync.Mutex
	lastProcessAt time.Time
	batchSize     int
	batchInterval time.Duration
	asyncProcess  bool
	closed        bool
}

func (h *BaseHandler) Init() error {
	conf := h.Config
	// matcher
	matcher, err := matcherOf(conf.Str("match", ""))
	if err != nil {
		return err
	}
	h.matcher = matcher

	// batchSize
	batchSize := conf.Int("batchSize", 20)
	h.batchSize = batchSize
	if h.batchSize > 0 {
		h.buff = make([]*Event, 0, h.batchSize)
	} else {
		h.buff = make([]*Event, 0, 100)
	}

	// formatter
	var formatterConf Config
	fv := conf.Interface("format", nil)
	if fv == nil {
		formatterConf = Config{"name": "json"}
	} else if fv1, ok := fv.(string); ok {
		formatterConf = Config{"name": fv1}
	} else if fv1, ok := fv.(map[string]interface{}); ok {
		formatterConf = Config(fv1)
	} else {
		return errors.New("Illegal format")
	}
	formatter, err := CompileFormatter(formatterConf)
	if err != nil {
		return err
	}
	h.Formatter = formatter

	// async
	h.asyncProcess = conf.Bool("async", false)

	// batchInterval
	batchIntervalSec := conf.Int("batchInterval", 60)
	if batchIntervalSec > 0 {
		h.batchInterval = time.Duration(batchIntervalSec) * time.Second
	}
	if h.batchInterval > 0 {
		go func(batchInterval time.Duration) {
			for {
				closed := false
				var processing []*Event
				lockDo(&h.buffMtx, func() {
					if time.Now().Sub(h.lastProcessAt) > h.batchInterval {
						processing = h.resetBuff()
					}
					if h.closed {
						closed = h.closed
					}
				})
				h.callProcess(processing)
				if closed {
					break
				}
				time.Sleep(h.batchInterval)
			}
		}(h.batchInterval)
	}
	return nil
}

func (h *BaseHandler) Close() error {
	lockDo(&h.buffMtx, func() {
		h.closed = true
	})
	return nil
}

func (h *BaseHandler) Match(e *Event) (bool, error) {
	return h.matcher(e)
}

func (h *BaseHandler) Append(e *Event) {
	var processing []*Event
	lockDo(&h.buffMtx, func() {
		h.buff = append(h.buff, e)
		proc := false
		if h.batchSize > 0 && len(h.buff) >= h.batchSize {
			proc = true
		}
		if proc {
			processing = h.resetBuff()
		}
	})
	h.callProcess(processing)
}

func (h *BaseHandler) Flush() {
	h.callProcess(h.resetBuff())
}

func (h *BaseHandler) doProcess(events []*Event) error {
	if h.Processor == nil {
		return errors.New("Nil processor")
	}
	return h.Processor(events)
}

func (h *BaseHandler) callProcess(events []*Event) {
	if len(events) == 0 {
		return
	}
	if h.asyncProcess {
		go safeDo(func() {
			err := h.doProcess(events)
			if err != nil {
				// TODO
			}
		})
	} else {
		safeDo(func() {
			err := h.doProcess(events)
			if err != nil {
				// TODO
			}
		})
	}
}

func (h *BaseHandler) resetBuff() []*Event {
	h.lastProcessAt = time.Now()
	if len(h.buff) == 0 {
		return nil
	}
	processing := h.buff
	h.buff = make([]*Event, 0, h.batchSize)
	return processing
}
