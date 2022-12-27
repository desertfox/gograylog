//GoGraylog is a simple client for interacting with a graylog instance.
package gograylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const SessionsPath string = "api/system/sessions"

var errMissingHost error = errors.New("client host is empty")

type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Host, token string
	query       query
	HttpClient  HTTPInterface
}

//
func (c *Client) Login(user, pass string) error {
	if c.Host == "" {
		return errMissingHost
	}

	lr := struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Host:     c.Host,
		Username: user,
		Password: pass,
	}

	data, err := json.Marshal(lr)
	if err != nil {
		return fmt.Errorf("error unable to encode loging request %w", err)
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%v/%v", c.Host, SessionsPath),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("error unable to construct request %w", err)
	}

	request.Header = defaultHeader()

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error making request %w", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body %w", err)
	}

	var respData map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("error unable to decode json data")
	}

	c.token = createAuthHeader(respData["session_id"] + ":session")

	return nil
}

//Execute query for a given stream id with specified fields and limit at a frequnecy.
func (c *Client) Search(search, streamid string, fields []string, limit, frequency int) ([]byte, error) {
	if c.token == "" {
		return nil, errors.New("no session found")
	}

	c.query = query{
		host:      c.Host,
		query:     search,
		streamid:  streamid,
		fields:    fields,
		limit:     limit,
		frequency: frequency,
	}

	return c.request(c.query)
}

//Generate GrayLog URL string between from-to
func (c Client) QueryFromTo(from, to time.Time) string {
	return c.query.interval(from, to)
}

func (c *Client) request(q query) ([]byte, error) {
	body, err := q.body()
	if err != nil {
		return nil, fmt.Errorf("unable to build body %w", err)
	}

	request, err := http.NewRequest("POST", q.url(), body)
	if err != nil {
		return nil, fmt.Errorf("unable to build request %w", err)
	}

	c.setHeaders(request)

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error submiting request %w", err)
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func defaultHeader() http.Header {
	h := &http.Header{}

	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", "GoGrayLog")

	return *h
}

func (c *Client) setHeaders(r *http.Request) {
	h := defaultHeader()

	h.Add("Authorization", c.token)
	h.Add("Accept", "text/csv")

	r.Header = h
}

func createAuthHeader(s string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
