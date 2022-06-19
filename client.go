package gograylog

import (
	"io"
	"net/http"
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
	request, _ := http.NewRequest("GET", q.URL(), q.BodyData())
	request.Close = true

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

	return body, nil
}

func defaultHeader() http.Header {
	h := http.Header{}

	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", "GoGrayLog 1")

	return h
}
