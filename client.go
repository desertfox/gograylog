package gograylog

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	DEBUG = os.Getenv("GOGRAYLOG_DEBUG")
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

func (c *Client) Execute(query, streamid string, frequency int) ([]byte, error) {
	return c.request(Query{
		Host:      c.session.loginRequest.Host,
		Query:     query,
		Streamid:  streamid,
		Frequency: frequency,
	})
}

func (c *Client) request(q Query) ([]byte, error) {
	if DEBUG != "" {
		fmt.Printf("query:%v\n", q)
	}

	request, _ := http.NewRequest("GET", q.URL(), q.BodyData())
	request.Close = true

	if DEBUG != "" {
		fmt.Printf("request body:%s\n", q.BodyData())
	}

	h := defaultHeader()
	h.Add("Authorization", c.session.authHeader())
	request.Header = h

	response, err := c.httpClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	if DEBUG != "" {
		fmt.Printf("response body:%s\n", body)
	}

	return body, nil
}

func defaultHeader() http.Header {
	h := http.Header{}

	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", "GoGrayLog 1")

	return h
}
