package conf

import "github.com/BurntSushi/toml"

// Config --
type Config struct {
	Project string `toml:"project"`
	Storage struct {
		Path string `toml:"path"`
	} `toml:"storage"`
	Source struct {
		Files []string `toml:"files"`
	}
}

// ReadConf --
func ReadConf(path string) (*Config, error) {
	conf := new(Config)
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
