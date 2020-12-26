package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleRequest(c *gin.Context) {
	endpoint := c.Param("endpoint")
	key := c.Param("key")
	reader := bufio.NewReader(c.Request.Body)

	err := config.processRequest(endpoint, key, c.Request.Header.Get("Content-Type"), reader)
	if err == nil {
		c.Status(http.StatusNoContent)
		return
	}

	errStr := err.Error()
	if len(errStr) == 3 {
		// Assume it's a status code.
		statusCode, err := strconv.Atoi(errStr)
		if err == nil {
			c.Status(statusCode)
			return
		}
		errStr = fmt.Sprintf("%s", err)
	}

	c.Data(500, gin.MIMEPlain, []byte(errStr))
}

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
