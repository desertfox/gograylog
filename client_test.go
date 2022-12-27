package gograylog

import (
	"errors"
	"net/http"
	"testing"
)

var (
	tUser, tPass = "X", "X"
)

type httpClientMock struct {
	response *http.Response
	error    error
}

func (h *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.response, h.error
}

func Test_LoginEmptyHost(t *testing.T) {
	c := Client{}

	got := c.Login(tUser, tPass)
	want := errMissingHost

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func Test_LoginRequest(t *testing.T) {
	c := Client{
		Host: "potato",
		HttpClient: &httpClientMock{
			response: nil,
			error:    errors.New("some error"),
		},
	}

	got := c.Login(tUser, tPass)
	if got == nil {
		t.Errorf("got %s", got)
	}
}
