package utils

import "flag"

func GetConfigPath() string {
	var fileConfig string
	flag.StringVar(&fileConfig, "config", "./cfg.json", "path to config file")
	flag.Parse()
	if fileConfig == "" {
		panic("No file directory")
	}
	return fileConfig
}
