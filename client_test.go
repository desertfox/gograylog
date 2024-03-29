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
	"time"
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
	testResponses := make(map[string]interface{})
	testResponses["session_id"] = "session id"
	testResponses["username"] = "duck"
	testResponses["valid_until"] = "2025-08-26T04:46:59.471+0000"
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
			expected: errMissingSession,
		},
		{
			description: "valid login",
			client: Client{
				Host: "potato",
				Session: &Session{
					Id:         "duck",
					ValidUntil: ValidUntil(time.Now()),
				},
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
			description: "missing Session",
			client: Client{
				Host:    "host",
				Session: &Session{},
			},
			query:         Query{},
			expectedError: errMissingSession,
			expectedData:  []byte{},
		},
		{
			description: "http error",
			client: Client{
				Host: "potato",
				Session: &Session{
					Id:         "duck",
					ValidUntil: ValidUntil(time.Now().Add(1 * time.Hour)),
				},
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
			description: "valid search with Session",
			client: Client{
				Host: "potato",
				Session: &Session{
					Id:         "duck",
					ValidUntil: ValidUntil(time.Now().Add(1 * time.Hour)),
				},
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

func Test_Streams(t *testing.T) {
	cases := []struct {
		description   string
		client        Client
		expectedError error
		expectedData  []byte
	}{
		{
			description:   "missing host",
			client:        Client{},
			expectedError: errMissingHost,
			expectedData:  []byte{},
		},
		{
			description: "missing Session",
			client: Client{
				Host:    "host",
				Session: &Session{},
			},
			expectedError: errMissingSession,
			expectedData:  []byte{},
		},
		{
			description: "http error",
			client: Client{
				Host: "potato",
				Session: &Session{
					Id:         "duck",
					ValidUntil: ValidUntil(time.Now().Add(1 * time.Hour)),
				},
				HttpClient: &httpClientMock{
					response: nil,
					error:    errTestHTTP,
				},
			},
			expectedError: errTestHTTP,
			expectedData:  []byte{},
		},
		{
			description: "valid streams with Session",
			client: Client{
				Host: "potato",
				Session: &Session{
					Id:         "duck",
					ValidUntil: ValidUntil(time.Now().Add(1 * time.Hour)),
				},
				HttpClient: &httpClientMock{
					response: &http.Response{
						Body: io.NopCloser(newBuf()),
					},
					error: nil,
				},
			},
			expectedError: nil,
			expectedData:  newBuf().Bytes(),
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			b, result := tt.client.Streams()
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
