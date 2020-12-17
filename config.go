package main

import (
	"encoding/json"
	"io/ioutil"
)

type Endpoint struct {
	Region       string `json:"region"`
	Bucket       string `json:"s3_bucket"`
	Prefix       string `json:"s3_prefix"`
	Queue        string `json:"sqs_queue"`
	Group        string `json:"sqs_group"`
	DelaySeconds int64  `json:"sqs_delay_seconds"`
}

type Config map[string]Endpoint

func LoadConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	res := Config{}
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
