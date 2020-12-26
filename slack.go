package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SlackPostData struct {
	Username  string `json:"username"`
	Channel   string `json:"channel"`
	Text      string `json:"text"`
	IconUrl   string `json:"icon_url"`
	IconEmoji string `json:"icon_emoji"`
}

func (config Endpoint) processSlackRequest(endpointName string, key string, contentType string, reader *bufio.Reader) error {
	postBody, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}

	postData := config.DefaultPostData

	switch contentType {
	case "application/json":
		if !config.AllowOverrides {
			return fmt.Errorf("403")
		}

		err = json.Unmarshal(postBody, &postData)
		if err != nil {
			return fmt.Errorf("failed to decode request body: %v", err)
		}

	default:
		message := string(postBody)
		if len(message) > 0 {
			postData.Text = message
		}
	}

	b, err := json.Marshal(postData)
	if err != nil {
		return fmt.Errorf("failed to encode Slack request body: %v", err)
	}

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Slack request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unsuccessful response from Slack: [%d] %s", resp.StatusCode, string(body))
	}

	return nil
}
