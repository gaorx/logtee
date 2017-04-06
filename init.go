package logtee

func init() {
	// formatter
	RegisterFormatterFactory("json", FormatterFactoryOf(jsonFormatter))
	RegisterFormatterFactory("kvl", FormatterFactoryOf(kvlFormatter))
	RegisterFormatterFactory("csv", newCsvFormatter)

	// handler
	RegisterHandlerFactory("scribe", NewScribeHandler)
}
