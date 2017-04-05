package logtee

type ScribeHandler struct {
	BaseHandler
}

func NewScribeHandler(name string, conf Config) Handler {
	return &ScribeHandler{
		BaseHandler: BaseHandler{
			Name:   name,
			Config: conf,
		},
	}
}

func (h *ScribeHandler) Init() error {
	err := h.BaseHandler.Init()
	if err != nil {
		return err
	}

	return nil
}

func (h *ScribeHandler) doProcess(events []*Event) error {
	println("****22", len(events))
	return nil
}
