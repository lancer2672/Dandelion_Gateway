package helper

import (
	"io"
	"net/http"

	"github.com/lancer2672/Dandelion_Gateway/internal/utils"
)

type ClientHttp struct {
	Client *http.Client
	Header http.Header
}

var HttpClient *ClientHttp

func newClientHttp(options map[string][]string) *ClientHttp {
	header := http.Header{}
	for key, values := range options {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return &ClientHttp{
		Client: &http.Client{},
		Header: header,
	}
}

func (c *ClientHttp) Do(req *http.Request) (*http.Response, error) {
	req.Header = c.Header
	return c.Client.Do(req)
}
func (c *ClientHttp) Get(url string) (r *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err

	}

	return c.Do(req)

}

func (c *ClientHttp) Post(url, contentType string, body io.Reader) (r *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)

}

func ConfigHttpClient(config utils.Config) error {
	options := make(map[string][]string)
	options["X-Gateway-Key"] = []string{config.GatewayApiKey}
	HttpClient = newClientHttp(options)

	return nil
}
