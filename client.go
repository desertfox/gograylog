// GoGraylog is a simple client for interacting with a graylog instance.
package gograylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	VERSION    string = "v1.6.5"
	acceptJSON string = "application/json"
	acceptCSV  string = "text/csv"
)

var (
	errMissingHost error             = errors.New("client host is empty")
	RouteMap       map[string]string = map[string]string{
		"sessions": "api/system/sessions",
		"messages": "api/views/search/messages",
		"streams":  "api/streams",
	}
)

type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientInterface interface {
	Login(string, string) error
	Search(QueryInterface) ([]byte, error)
}

// Graylog SDK client
type Client struct {
	Host       string
	Session    *Session
	HttpClient HTTPInterface
}

func (c *Client) Login(user, pass string) error {
	if c.Host == "" {
		return errMissingHost
	}

	data, err := json.Marshal(struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Host:     c.Host,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return fmt.Errorf("error unable to encode login request %w", err)
	}

	body, err := c.httpRequest(http.MethodPost, "sessions", bytes.NewBuffer(data), acceptJSON, false)
	if err != nil {
		return err
	}

	var session Session
	err = json.Unmarshal(body, &session)
	if err != nil {
		return fmt.Errorf("error unable to decode json data, %v %s", err, string(body))
	}

	if err = session.IsValid(); err != nil {
		return err
	}

	c.Session = &session

	return nil
}

// Execute Graylog search using GoGrayLog Query type
func (c *Client) Search(q QueryInterface) ([]byte, error) {
	if c.Host == "" {
		return nil, errMissingHost
	}

	if err := c.Session.IsValid(); err != nil {
		return nil, err
	}

	body, err := q.JSON()
	if err != nil {
		return nil, err
	}

	return c.httpRequest(http.MethodPost, "messages", bytes.NewReader(body), acceptCSV, true)
}

// Requests the Streams for the configured Client.Host graylog instance.
func (c *Client) Streams() ([]byte, error) {
	if c.Host == "" {
		return nil, errMissingHost
	}

	if err := c.Session.IsValid(); err != nil {
		return nil, err
	}

	return c.httpRequest(http.MethodGet, "streams", nil, acceptJSON, true)
}

func (c *Client) httpRequest(method, route string, body io.Reader, accept string, sendAuth bool) ([]byte, error) {
	if _, ok := RouteMap[route]; !ok {
		return nil, fmt.Errorf("route not found, %s", route)
	}
	r, err := http.NewRequest(
		method,
		fmt.Sprintf("%v/%v", c.Host, RouteMap[route]),
		body,
	)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("X-Requested-By", fmt.Sprintf("GoGrayLog %s", VERSION))
	r.Header.Add("Accept", accept)
	if sendAuth {
		r.Header.Add("Authorization", createAuthHeader(c.Session.Id+":session"))
	}

	response, err := c.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body %w", err)
	}

	return buf.Bytes(), nil
}

func createAuthHeader(s string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
