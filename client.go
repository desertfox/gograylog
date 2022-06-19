package main

import (
	"io"
	"net/http"
)

var (
	VERSION = "0.1"
)

type Client struct {
	session    *session
	httpClient *http.Client
}

func New(host, user, pass string) *Client {
	var httpClient *http.Client = &http.Client{}
	return &Client{
		session:    newSession(host, user, pass, httpClient),
		httpClient: httpClient,
	}
}

func (c *Client) SearchStream(q Query) ([]byte, error) {
	return c.request(q)
}

func (c *Client) request(q Query) ([]byte, error) {
	request, _ := http.NewRequest("GET", q.URL(), q.BodyData())
	request.Close = true

	request.Header.Set("Authorization", c.session.authHeader())
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("X-Requested-By", "GoGrayLog "+VERSION)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
