package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	remoteData := []byte(`{"result": "triggered successfully"}`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/some/path", r.URL.String())

		responseBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "Expected no error reading request body inside test")
		assert.Equal(t, "payload", string(responseBody), "Expected request body to match")

		w.Write(remoteData)
	}))
	defer ts.Close()

	req := Request{
		BaseURL: ts.URL,
		Route:   "/some/path",
		Method:  "POST",
		Body:    bytes.NewBufferString("payload"),
	}
	data, err := NewRemoteFetcher().Fetch(req)

	require.NoError(t, err, "Expected no error for successful response")
	assert.Equal(t, `{"result": "triggered successfully"}`, string(data), "Expected response body to match")
}

func TestFetchWithAuthorizationHeader(t *testing.T) {
	remoteData := []byte(`{"group": "My Group"}`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer eyj0BlahBlah", r.Header.Get("authorization"), "Expected request to have custom header")

		w.Write(remoteData)
	}))
	defer ts.Close()

	req := Request{
		BaseURL: ts.URL,
		Route:   "/some/path",
		Method:  "POST",
		Body:    bytes.NewBufferString("payload"),
		Headers: map[string]string{"authorization": "Bearer eyj0BlahBlah"},
	}
	_, err := NewRemoteFetcher().Fetch(req)

	require.NoError(t, err, "Expected no error for successful response")
}

func TestFetchWithCustomHeaders(t *testing.T) {
	remoteData := []byte(`{"result": "triggered successfully"}`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "FooValue", r.Header.Get("Foo"), "Expected request to have custom header")
		assert.Equal(t, "BarValue", r.Header.Get("Bar"), "Expected request to have custom header")

		w.Write(remoteData)
	}))
	defer ts.Close()

	req := Request{
		BaseURL: ts.URL,
		Route:   "/some/path",
		Body:    bytes.NewBufferString("payload"),
		Headers: map[string]string{"Foo": "FooValue", "Bar": "BarValue"},
	}
	_, err := NewRemoteFetcher().Fetch(req)

	require.NoError(t, err, "Expected no error for successful response")
}

func TestFetchWhenRequestCreationFails(t *testing.T) {
	req := Request{
		BaseURL: "",
		Method:  "GET",
		Route:   "/some/path",
		Body:    bytes.NewBufferString("payload"),
	}
	data, err := NewRemoteFetcher().Fetch(req)

	require.Error(t, err, "Expected error for when request cannot be send")
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
	assert.Nil(t, data, "Expected no invalid response body")
}

func TestFetchWhenSendingRequestFails(t *testing.T) {
	req := Request{
		BaseURL: "http://<>",
		Route:   "/some/path",
		Method:  "POST",
		Body:    bytes.NewBufferString("payload"),
	}
	data, err := NewRemoteFetcher().Fetch(req)

	require.Error(t, err, "Expected error for when request cannot be send")
	assert.Contains(t, err.Error(), "POST http://<>/some/path: server: ")
	assert.Nil(t, data, "Expected no invalid response body")
}

func TestFetchWhenHTTPStatusNotSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "database not reachable", http.StatusInternalServerError)
	}))
	defer ts.Close()

	req := Request{
		BaseURL: ts.URL,
		Route:   "/some/path",
		Method:  "POST",
		Body:    bytes.NewBufferString("payload"),
	}
		data, err := NewRemoteFetcher().Fetch(req)

	url := fmt.Sprintf("%s%s", req.BaseURL, req.Route)

	require.Error(t, err, "Expected error for when server returns non-success HTTP status code")
	assert.EqualError(t, err, fmt.Sprintf("POST %s: HTTP status code was 500, body: database not reachable\n", url))
	assert.Nil(t, data, "Expected no invalid response body")
}

type badReader struct{}

func (br badReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestParseResponseForBadBody(t *testing.T) {
	badBody := ioutil.NopCloser(badReader{})

	responseBody, err := parseResponse(&http.Response{Body: badBody, StatusCode: http.StatusOK})

	require.Error(t, err, "Expected error on read")
	assert.Contains(t, err.Error(), "read error")
	assert.Nil(t, responseBody)
}
