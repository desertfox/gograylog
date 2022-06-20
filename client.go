package gograylog

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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
		fmt.Printf("url:%v\n", q.URL())
		fmt.Printf("query:%v\n", q)
	}

	request, _ := http.NewRequest("POST", q.URL(), q.BodyData())

	h := defaultHeader()
	h.Add("Authorization", c.session.authHeader())
	h.Add("Accept", "text/csv")
	request.Header = h

	if DEBUG != "" {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Printf("request: %q\n", dump)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	if DEBUG != "" {
		dump, _ := httputil.DumpResponse(response, true)

		fmt.Printf("response: %q\n", dump)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func defaultHeader() http.Header {
	h := http.Header{}

	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", "GoGrayLog")

	return h
}
