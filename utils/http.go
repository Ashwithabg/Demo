package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	BaseURL string
	Route   string
	Method  string
	Body    io.Reader
	Headers map[string]string
}

type Fetcher interface {
	Fetch(req Request) (responseBody []byte, err error)
}

type RemoteFetcher struct {
	client *http.Client
}

func NewRemoteFetcher() *RemoteFetcher {
	return &RemoteFetcher{client: &http.Client{}}
}

func (fetcher RemoteFetcher) Fetch(req Request) ([]byte, error) {
	method := req.Method
	url := fmt.Sprintf("%s%s", req.BaseURL, req.Route)

	httpReq, err := http.NewRequest(method, url, req.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	res, err := fetcher.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s %s: server: %s", method, url, err)
	}
	defer res.Body.Close()

	body, err := parseResponse(res)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", method, url, err)
	}

	return body, err
}

func parseResponse(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read HTTP response body: %s", err)
	}

	if !IsStatusCodeValid(res) {
		return body, fmt.Errorf("HTTP status code was %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}

func IsStatusCodeValid(resp *http.Response) bool {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false
	}
	return true
}
