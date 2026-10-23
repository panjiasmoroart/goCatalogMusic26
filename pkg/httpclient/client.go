package httpclient

import (
	"net/http"
)

//go:generate mockgen -source=client.go -destination=client_mock.go -package=httpclient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Client HTTPClient
}

func NewClient(client HTTPClient) *Client {
	return &Client{
		Client: client,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// fmt.Printf("http client: %+v\n", c.Client)

	return c.Client.Do(req)
}
