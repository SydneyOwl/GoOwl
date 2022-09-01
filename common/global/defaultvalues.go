package global

const (
	Version = "0.3.3"
	Status  = "release"
)

var (
	LoggingMethod     = 1
	OS                = "linux"
	ErrorBuildStatus  = "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"74\" height=\"20\" role=\"img\" aria-label=\"Build: Error\"><title>Build: Error</title><linearGradient id=\"s\" x2=\"0\" y2=\"100%\"><stop offset=\"0\" stop-color=\"#bbb\" stop-opacity=\".1\"/><stop offset=\"1\" stop-opacity=\".1\"/></linearGradient><clipPath id=\"r\"><rect width=\"74\" height=\"20\" rx=\"3\" fill=\"#fff\"/></clipPath><g clip-path=\"url(#r)\"><rect width=\"37\" height=\"20\" fill=\"#555\"/><rect x=\"37\" width=\"37\" height=\"20\" fill=\"#e05d44\"/><rect width=\"74\" height=\"20\" fill=\"url(#s)\"/></g><g fill=\"#fff\" text-anchor=\"middle\" font-family=\"Verdana,Geneva,DejaVu Sans,sans-serif\" text-rendering=\"geometricPrecision\" font-size=\"110\"><text aria-hidden=\"true\" x=\"195\" y=\"150\" fill=\"#010101\" fill-opacity=\".3\" transform=\"scale(.1)\" textLength=\"270\">Build</text><text x=\"195\" y=\"140\" transform=\"scale(.1)\" fill=\"#fff\" textLength=\"270\">Build</text><text aria-hidden=\"true\" x=\"545\" y=\"150\" fill=\"#010101\" fill-opacity=\".3\" transform=\"scale(.1)\" textLength=\"270\">Error</text><text x=\"545\" y=\"140\" transform=\"scale(.1)\" fill=\"#fff\" textLength=\"270\">Error</text></g></svg>"
	Sqlite3DBPosition = "./GoOwl.db"
	SqlDebug          = false
)
