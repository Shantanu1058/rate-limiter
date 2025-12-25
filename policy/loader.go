package policy

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func LoadPolicyConfig(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
