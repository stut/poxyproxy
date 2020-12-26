package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// This contains the top-level configuration for all types.
type Endpoint struct {
	Type string `json:"type"`

	// S3
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
	Prefix   string `json:"prefix"`
	SQS      SQS    `json:"sqs"`

	// Slack
	Url             string        `json:"url"`
	AllowOverrides  bool          `json:"allow_overrides"`
	DefaultPostData SlackPostData `json:"defaults"`
}

type Endpoints map[string]Endpoint

func LoadEndpointsConfig(filename string) (*Endpoints, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := Endpoints{}
	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, err
	}

	configValid := true
	for endpointName, endpoint := range config {
		switch endpoint.Type {
		case "s3":
		case "slack":
			continue
		default:
			os.Stderr.WriteString(fmt.Sprintf("Unknown type \"%s\" configured for \"%s\"", endpoint.Type, endpointName))
			configValid = false
		}
	}

	if !configValid {
		return nil, fmt.Errorf("endpoints configuration is invalid, cannot continue!")
	}

	return &config, nil
}
