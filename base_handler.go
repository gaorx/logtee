package logtee

import (
	"sync"
	"time"
)

type BaseHandler struct {
	Name          string
	Config        Config
	matcher       matcher
	data          map[string]interface{}
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
	batchSize := conf.Int("batch_size", 20)
	h.batchSize = batchSize
	if h.batchSize > 0 {
		h.buff = make([]*Event, 0, h.batchSize)
	} else {
		h.buff = make([]*Event, 0, 100)
	}

	// async process
	h.asyncProcess = conf.Bool("async_process", false)

	// batchInterval
	batchIntervalSec := conf.Int("batch_interval", 60)
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

func (h *BaseHandler) doProcess(events []*Event) error {
	println("****11", len(events))
	return nil
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

func (h *BaseHandler) Set(k string, v interface{}) {
	if h.data == nil {
		h.data = map[string]interface{}{}
	}
	h.data[k] = v
}

func (h *BaseHandler) Get(k string, def interface{}) interface{} {
	v, ok := h.data[k]
	if !ok {
		return def
	}
	return v
}
