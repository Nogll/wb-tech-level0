package config

type AppConfig struct {
	DbUrl         string
	Brokers       []string
	Topic         string
	GroupId       string
	StaticDir     string
	IndexTemplate string
	Addr          string
	PreloadLimit  int32
}

var defaultConfig = AppConfig{
	DbUrl:         "postgres://user:password@localhost:5432/db",
	Brokers:       []string{"localhost:29092"},
	Topic:         "orders",
	GroupId:       "go-reader",
	StaticDir:     "static/",
	IndexTemplate: "templates/index.html",
	Addr:          ":8000",
	PreloadLimit:  2,
}

func LoadConfig() AppConfig {
	return defaultConfig
}
