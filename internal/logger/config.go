package logger

type Config struct {
	Targets string
	Levels  map[string]string
}

func DefaultConfig() *Config {
	conf := &Config{
		Levels:  make(map[string]string),
		Targets: "console",
	}

	conf.Levels["default"] = "info"
	conf.Levels["_transport"] = "info"
	conf.Levels["_protocol"] = "info"
	conf.Levels["_pool"] = "warn"

	return conf
}

// BasicCheck performs basic checks on the configuration.
func (*Config) BasicCheck() error {
	return nil
}
