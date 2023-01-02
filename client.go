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
)

const (
	SessionsPath string = "api/system/sessions"
	VERSION      string = "v0.1.0"
)

var (
	errMissingHost      error = errors.New("client host is empty")
	errMissingToken     error = errors.New("no token found on client")
	errMissingSessionID error = errors.New("response is missing session_id key")
)

type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Host, token string
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
		return fmt.Errorf("error unable to encode login request %w", err)
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%v/%v", c.Host, SessionsPath),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("error unable to construct request %w", err)
	}

	h := http.Header{}
	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", VERSION)
	request.Header = h

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body %w", err)
	}

	var respData map[string]string
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return fmt.Errorf("error unable to decode json data, %s", respData)
	}

	if _, ok := respData["session_id"]; !ok {
		return errMissingSessionID
	}

	c.token = createAuthHeader(respData["session_id"] + ":session")

	return nil
}

//Execute query for a given stream id with specified fields and limit at a frequnecy.
func (c *Client) Search(q Query) ([]byte, error) {
	if c.token == "" {
		return nil, errMissingToken
	}

	body, err := q.JSON()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", q.endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	h := http.Header{}
	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", VERSION)
	h.Add("Authorization", c.token)
	h.Add("Accept", "text/csv")
	request.Header = h

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func createAuthHeader(s string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
