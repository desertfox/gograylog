package gograylog

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
						Body: io.NopCloser(bytes.NewReader(testLoginResponse)),
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
		description   string
		client        Client
		query         Query
		expectedError error
		expectedData  []byte
	}{
		{
			description:   "missing host",
			client:        Client{},
			query:         Query{},
			expectedError: errMissingHost,
			expectedData:  []byte{},
		},
		{
			description: "missing token",
			client: Client{
				Host: "host",
			},
			query:         Query{},
			expectedError: errMissingAuth,
			expectedData:  []byte{},
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
				QueryString: "error",
				StreamID:    "somehash",
				Frequency:   15,
			},
			expectedError: errTestHTTP,
			expectedData:  []byte{},
		},
		{
			description: "valid search with token",
			client: Client{
				Host:  "potato",
				token: "sometoken",
				HttpClient: &httpClientMock{
					response: &http.Response{
						Body: io.NopCloser(newBuf()),
					},
					error: nil,
				},
			},
			query: Query{
				QueryString: "error",
				StreamID:    "somehash",
				Frequency:   15,
			},
			expectedError: nil,
			expectedData:  newBuf().Bytes(),
		},
		{
			description: "valid search with basicauth",
			client: Client{
				Host:     "potato",
				Username: "username",
				Password: "password",
				HttpClient: &httpClientMock{
					response: &http.Response{
						Body: io.NopCloser(newBuf()),
					},
					error: nil,
				},
			},
			query: Query{
				QueryString: "error",
				StreamID:    "somehash",
				Frequency:   15,
			},
			expectedError: nil,
			expectedData:  newBuf().Bytes(),
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			b, result := tt.client.Search(tt.query)
			if result != tt.expectedError {
				t.Errorf("expected %d, but got %d", tt.expectedError, result)
			}
			if string(b) != string(tt.expectedData) {
				t.Errorf("expected %s, but got %s", tt.expectedData, b)
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

func newBuf() *bytes.Buffer {
	var testSearchResponse *bytes.Buffer = bytes.NewBuffer([]byte{})
	records := [][]string{
		{"message"},
		{"m1 error"},
		{"m2 error"},
	}

	w := csv.NewWriter(testSearchResponse)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	return testSearchResponse
}
