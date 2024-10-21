package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readConfig(filePath string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return config, nil
}
