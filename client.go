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
	VERSION      string = "v1.4.0"
)

var (
	errMissingHost      error = errors.New("client host is empty")
	errMissingAuth      error = errors.New("no auth found on client")
	errMissingSessionID error = errors.New("response is missing session_id key")
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
	Token      string
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

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%v/%v", c.Host, SessionsPath),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("error unable to construct request %w", err)
	}

	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("X-Requested-By", fmt.Sprintf("GoGrayLog %s", VERSION))
	request.Header.Add("Accept", "application/json")

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

	c.Token = respData["session_id"]

	return nil
}

// Execute Graylog search using GoGrayLog Query type
func (c *Client) Search(q QueryInterface) ([]byte, error) {
	if c.Host == "" {
		return nil, errMissingHost
	}

	if c.Token == "" {
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

	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("X-Requested-By", fmt.Sprintf("GoGrayLog %s", VERSION))
	request.Header.Add("Accept", "text/csv")
	request.Header.Add("Authorization", createAuthHeader(c.Token+":session"))

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body %w", err)
	}

	return data, nil
}

func createAuthHeader(s string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
