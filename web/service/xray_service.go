package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type InboundConfig struct {
	Listen   string                 `json:"listen"`
	Port     int                    `json:"port"`
	Protocol string                 `json:"protocol"`
	Settings map[string]interface{} `json:"settings"`
	Tag      string                 `json:"tag"`
}

func WriteInboundConfigToFile(config *InboundConfig, filePath string) error {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal InboundConfig to JSON: %v", err)
	}

	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write InboundConfig to file: %v", err)
	}

	return nil
}
