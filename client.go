package jandi

// Jandi Webhook Library for Golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Incoming is a payload for an incoming webhook
type incoming struct {
	Body         string        `json:"body"`         // markdown supported
	ConnectColor string        `json:"connectColor"` // #RRGGBB format
	ConnectInfo  []ConnectInfo `json:"connectInfo"`
}

// ConnectInfo is the information of `Jandi connect`
type ConnectInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl,omitempty"`
}

// Constants for webhook
const (
	headerAccept      = "application/vnd.tosslab.jandi-v2+json"
	headerContentType = "application/json"
)

// IncomingClient is a client for sending incoming webhooks
type IncomingClient struct {
	webhookURL string
	verbose    bool
}

// NewIncomingClient creates a new IncomingClient
func NewIncomingClient(webhookURL string) *IncomingClient {
	return &IncomingClient{
		webhookURL: webhookURL,
		verbose:    false,
	}
}

// ConnectInfoFrom generates an array of one ConnectInfo filled with given title, description, and image URL
func ConnectInfoFrom(title, description, imageURL string) []ConnectInfo {
	return []ConnectInfo{
		ConnectInfo{
			Title:       title,
			Description: description,
			ImageURL:    imageURL,
		},
	}
}

// ConnectInfoNone returns an empty array of ConnectInfo
func ConnectInfoNone() []ConnectInfo {
	return []ConnectInfo{}
}

// SetVerbose sets if verbose error messages are shown or not
func (c *IncomingClient) SetVerbose(verbose bool) {
	c.verbose = verbose
}

// SendIncoming sends an incoming webhook
func (c *IncomingClient) SendIncoming(body, color string, infos []ConnectInfo) (result string, err error) {
	payload := incoming{
		Body:         body,
		ConnectColor: color,
		ConnectInfo:  infos,
	}

	var data, txt []byte
	if data, err = json.Marshal(payload); err == nil {
		reader := bytes.NewReader(data)
		var req *http.Request
		if req, err = http.NewRequest("POST", c.webhookURL, reader); err == nil {
			req.Header.Add("Accept", headerAccept)
			req.Header.Add("Content-Type", headerContentType)

			var resp *http.Response
			client := &http.Client{}
			resp, err = client.Do(req)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err == nil {
				txt, _ = ioutil.ReadAll(resp.Body)

				if c.verbose {
					log.Printf("<<< %+v", req)
				}

				if resp.StatusCode != 200 {
					if len(txt) > 0 {
						err = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(txt))
					} else {
						err = fmt.Errorf("HTTP %d", resp.StatusCode)
					}

					if c.verbose {
						log.Printf(">>> %s", err)
					}
				} else {
					return string(txt), nil
				}
			}
		}
	}

	return string(txt), err
}
