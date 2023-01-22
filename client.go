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
	//Endpoint to attempt login to
	SessionsPath string = "api/system/sessions"
	MessagesPath string = "api/views/search/messages"
	VERSION      string = "v1.2.0"
)

var (
	errMissingHost      error = errors.New("client host is empty")
	errMissingAuth      error = errors.New("no auth found on client")
	errMissingSessionID error = errors.New("response is missing session_id key")
)

type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Graylog SDK client
type Client struct {
	Host       string
	Username   string
	Password   string
	HttpClient HTTPInterface
	token      string
}

// Method to execute login request to the configured Client.Host
// if this method is used to login username and password will be discarded
// and the session token will be used for request authorization
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

	body, err := io.ReadAll(response.Body)
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

// Execute Graylog search using GoGrayLog Query type
func (c *Client) Search(q Query) ([]byte, error) {
	if c.Host == "" {
		return nil, errMissingHost
	}

	if c.token == "" && (c.Username == "" || c.Password == "") {
		return nil, errMissingAuth
	}

	body, err := q.JSON()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%v/%v", c.Host, MessagesPath),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	h := http.Header{}
	h.Add("Content-Type", "application/json; charset=UTF-8")
	h.Add("X-Requested-By", VERSION)
	h.Add("Accept", "text/csv")

	if c.token != "" {
		h.Add("Authorization", c.token)
	} else {
		request.SetBasicAuth(c.Username, c.Password)
	}

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
