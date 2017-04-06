package logtee

type ScribeHandler struct {
	BaseHandler
}

func NewScribeHandler(name string, conf Config) Handler {
	h := &ScribeHandler{
		BaseHandler: BaseHandler{
			Name:   name,
			Config: conf,
		},
	}
	h.Processor = func(events []*Event) error {
		return h.process(events)
	}
	return h
}

func (h *ScribeHandler) Init() error {
	err := h.BaseHandler.Init()
	if err != nil {
		return err
	}
	return nil
}

func (h *ScribeHandler) process(events []*Event) error {
	for _, e := range events {
		b, _ := h.Formatter(e)
		println("**", string(b))
	}
	return nil
}
