//GoGraylog is a simple client for interacting with a graylog instance.
package gograylog

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	query      query
	session    *session
	httpClient *http.Client
}

//Constructor
func New(host, user, pass string) *Client {
	var httpClient *http.Client = &http.Client{}

	return &Client{
		session:    newSession(host, user, pass, httpClient),
		httpClient: httpClient,
	}
}

//Execute query for a given stream id with specified fields and limit at a frequnecy.
func (c *Client) Execute(search, streamid string, fields []string, limit, frequency int) ([]byte, error) {
	c.session.login()

	c.query = query{
		host:      c.Host(),
		query:     search,
		streamid:  streamid,
		fields:    fields,
		limit:     limit,
		frequency: frequency,
	}

	return c.request(c.query)
}

//Return session Host
func (c Client) Host() string {
	return c.session.host()
}

//Generate GrayLog URL string between from-to
func (c Client) BuildURL(from, to time.Time) string {
	return c.query.string(from, to)
}

func (c *Client) request(q query) ([]byte, error) {
	request, _ := http.NewRequest("POST", q.url(), q.data())

	c.setHeaders(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func defaultHeader() *http.Header {
	h := &http.Header{}

	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", "GoGrayLog")

	return h
}

func (c *Client) setHeaders(r *http.Request) {
	h := defaultHeader()

	c.session.setAuthHeader(h)

	h.Add("Accept", "text/csv")

	r.Header = *h
}
