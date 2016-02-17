package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Codes struct {
	Weather map[string]string
}

var(
	WeatherCodes = &Codes{}
)

func init() {
	configYaml, readErr := ioutil.ReadFile("codes.yaml")
	if readErr != nil {
		return
	}
	yamlErr := yaml.Unmarshal(configYaml, WeatherCodes)
	if yamlErr != nil {
		return
	}
}
