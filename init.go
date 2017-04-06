package logtee

func init() {
	// formatter
	RegisterFormatterFactory("json", FormatterFactoryOf(jsonFormatter))
	RegisterFormatterFactory("csv", newCsvFormatter)

	// handler
	RegisterHandlerFactory("scribe", NewScribeHandler)
}
