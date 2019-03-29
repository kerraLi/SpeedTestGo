package util

import (
	"github.com/kylelemons/go-gypsy/yaml"
)

var file, _ = yaml.ReadFile("src/main/config/config.yml")

func GetConfig(value string) string {
	val, err := file.Get(value)
	if err != nil {
		return ""
	}
	return val
}
