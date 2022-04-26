package conf

import "github.com/BurntSushi/toml"

// Config --
type Config struct {
	Project string `toml:"project"`
	Version string `toml:"version"`
	Storage struct {
		Path string `toml:"path"`
	} `toml:"storage"`
	Source struct {
		Files []string `toml:"files"`
	} `toml:"source"`
	Merge struct {
		ChannelSize uint64 `toml:"channel_size"`
	} `toml:"merge"`
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
