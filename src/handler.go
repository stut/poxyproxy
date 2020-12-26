package main

import (
	"bufio"
	"fmt"
)

func (endpoints Endpoints) processRequest(endpointName string, key string, contentType string, reader *bufio.Reader) (err error) {
	endpoint, gotEndpoint := endpoints[endpointName]
	if !gotEndpoint {
		return fmt.Errorf("404")
	}

	switch endpoint.Type {
	case "s3":
		return endpoint.processS3Request(endpointName, key, contentType, reader)
	case "slack":
		return endpoint.processSlackRequest(endpointName, key, contentType, reader)
	default:
		return fmt.Errorf("endpoint configuration is invalid")
	}
}
