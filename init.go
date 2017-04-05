package logtee

func init() {
	RegisterHandlerFactory("scribe", NewScribeHandler)
}
