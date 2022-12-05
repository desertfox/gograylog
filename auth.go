package gograylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const sessionsPath string = "api/system/sessions"

var sessionInstanceMap = make(map[string]*session)

type session struct {
	basicAuth    string
	updated      time.Time
	loginRequest *loginRequest
	httpClient   *http.Client
}

type loginRequest struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func newSession(host, user, pass string, httpClient *http.Client) *session {
	if _, exists := sessionInstanceMap[host]; !exists {
		sessionInstanceMap[host] = &session{
			"",
			time.Now(),
			&loginRequest{host, user, pass},
			httpClient,
		}
	}

	return sessionInstanceMap[host]
}

func (s *session) host() string {
	return s.loginRequest.Host
}

func (s *session) login() {
	if s.basicAuth == "" {
		s.buildBasicAuth()
	}
}

func (s *session) buildBasicAuth() {
	sessionId, err := s.loginRequest.execute(s.httpClient)
	if err != nil {
		panic(err.Error())
	}
	s.basicAuth = createAuthHeader(sessionId + ":session")

	s.updated = time.Now()
}

func (s *session) setAuthHeader(h *http.Header) {
	h.Add("Authorization", s.basicAuth)
}

func (lr loginRequest) execute(httpClient *http.Client) (string, error) {
	request := lr.create()

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data["session_id"], nil
}

func (lr loginRequest) create() *http.Request {
	jsonData, err := json.Marshal(lr)
	if err != nil {
		panic(err)
	}
	request, _ := http.NewRequest("POST", fmt.Sprintf("%v/%v", lr.Host, sessionsPath), bytes.NewBuffer(jsonData))
	request.Header = *defaultHeader()

	return request
}

func createAuthHeader(s string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
