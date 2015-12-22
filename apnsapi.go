package apnsapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DevelopmentServer = "https://api.development.push.apple.com"
	ProductionServer  = "https://api.push.apple.com"
)

func NewClient(host string, client *http.Client) *Client {
	return &Client{
		host:   host,
		client: client,
	}
}

type Client struct {
	host   string
	client *http.Client
}

type Header struct {
	ApnsID         string
	ApnsExpiration string
	ApnsPriority   string
	ApnsTopic      string
}

type Response struct {
	ApnsID     string
	StatusCode int
}

type ErrorResponse struct {
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"` //TODO is string ?
}

func (e *ErrorResponse) Error() string {
	return e.Reason
}

func (c *Client) Do(token string, header *Header, payload []byte) (*Response, error) {
	req, err := c.NewRquest(token, header, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r := &Response{
		ApnsID:     res.Header.Get("apns-id"),
		StatusCode: res.StatusCode,
	}

	if res.StatusCode == http.StatusOK {
	} else {
		var er ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
			return r, err
		}
		return r, &er
	}

	return r, nil
}

func (c *Client) NewRquest(token string, header *Header, payload io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/3/device/%s", c.host, token)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if header != nil {
		if header.ApnsID != "" {
			req.Header.Set("apns-id", header.ApnsID)
		}
		if header.ApnsExpiration != "" {
			req.Header.Set("apns-expiration", header.ApnsExpiration)
		}
		if header.ApnsPriority != "" {
			req.Header.Set("apns-priority", header.ApnsPriority)
		}
		if header.ApnsTopic != "" {
			req.Header.Set("apns-topic", header.ApnsTopic)
		}
	}

	return req, err
}
