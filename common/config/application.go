package config

type ApplicationSettings struct {
	Mode string
	Host string
	Name string
	Port int
}

var ApplicationConfig = new(ApplicationSettings)
