package config

type LogConfig struct {
	Output string `default:"stdout" validate:"oneof=stdout file"`
	Level  string `default:"debug" validate:"oneof=trace debug info"`
	Path   string `default:"./log/" validate:"dirpath"`
}
