package config

import (
	"github.com/Unknwon/goconfig"
)

func LoadConfig() (*goconfig.ConfigFile, error) {
	cfg, err := goconfig.LoadConfigFile("./conf.ini")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
