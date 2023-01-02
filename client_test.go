package gograylog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
)

var (
	tUser, tPass            = "X", "X"
	errTestHTTP       error = errors.New("some error")
	testLoginResponse []byte
)

type httpClientMock struct {
	response *http.Response
	error    error
}

func (h *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.response, h.error
}

func init() {
	testResponses := make(map[string]string)
	testResponses["session_id"] = "session id"
	testLoginResponse, _ = json.Marshal(testResponses)
}

func Test_Login(t *testing.T) {
	cases := []struct {
		description        string
		client             Client
		username, password string
		expected           error
	}{
		{
			description: "missing host",
			client:      Client{},
			username:    "nohost",
			password:    "secret",
			expected:    errMissingHost,
		},
		{
			description: "http error",
			client: Client{
				Host: "potato",
				HttpClient: &httpClientMock{
					response: nil,
					error:    errTestHTTP,
				},
			},
			username: "httperror",
			password: "secret",
			expected: errTestHTTP,
		},
		{
			description: "missing session id",
			client: Client{
				Host: "potato",
				HttpClient: &httpClientMock{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(`{"something": "potato"}`))),
					},
					error: nil,
				},
			},
			username: "missingsessionid",
			password: "secret",
			expected: errMissingSessionID,
		},
		{
			description: "valid login",
			client: Client{
				Host: "potato",
				HttpClient: &httpClientMock{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(`{"session_id": "SESSIONID"}`))),
					},
					error: nil,
				},
			},
			username: "validlogin",
			password: "secret",
			expected: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			result := tt.client.Login(tt.username, tt.password)
			if result != tt.expected {
				t.Errorf("expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func Test_Search(t *testing.T) {
	cases := []struct {
		description string
		client      Client
		query       Query
		expected    error
	}{
		{
			description: "missing token",
			client:      Client{},
			query:       Query{},
			expected:    errMissingToken,
		},
		{
			description: "http error",
			client: Client{
				Host:  "potato",
				token: "faketoken",
				HttpClient: &httpClientMock{
					response: nil,
					error:    errTestHTTP,
				},
			},
			query: Query{
				Host:        "https://desertfox.dev",
				QueryString: "error",
				StreamID:    "somehash",
				Frequency:   15,
			},
			expected: errTestHTTP,
		},
		{
			description: "valid search",
			client: Client{
				Host:  "potato",
				token: "sometoken",
				HttpClient: &httpClientMock{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader(testLoginResponse)),
					},
					error: nil,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			_, result := tt.client.Search(tt.query)
			if result != tt.expected {
				t.Errorf("expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func ExampleClient() {
	c := Client{}
	err := c.Login(tUser, tPass)
	fmt.Println(err)
	//Output ds
}
