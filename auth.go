package gograylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"
)

const sessionsPath string = "api/system/sessions"

var (
	lock               = &sync.Mutex{}
	sessionInstanceMap = make(map[string]*session)
)

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
	lock.Lock()
	defer lock.Unlock()

	if _, exists := sessionInstanceMap[host]; !exists {
		sessionInstanceMap[host] = &session{"", time.Now(), &loginRequest{host, user, pass}, httpClient}

		sessionInstanceMap[host].buildBasicAuth()
	}

	return sessionInstanceMap[host]
}

func (s *session) buildBasicAuth() {
	if os.Getenv("GOGRAYLOG_UNPW") == "1" {
		s.basicAuth = createAuthHeader(s.loginRequest.Username + ":" + s.loginRequest.Password)
	} else {
		sessionId, err := s.loginRequest.execute(s.httpClient)
		if err != nil {
			panic(err.Error())
		}
		s.basicAuth = createAuthHeader(sessionId + ":session")
	}
	s.updated = time.Now()
}

func (s *session) authHeader() string {
	if s.basicAuth == "" {
		s.buildBasicAuth()
	}
	return s.basicAuth
}

func (lr loginRequest) execute(httpClient *http.Client) (string, error) {
	jsonData, err := json.Marshal(lr)
	if err != nil {
		return "", err
	}

	request, _ := http.NewRequest("POST", fmt.Sprintf("%v/%v", lr.Host, sessionsPath), bytes.NewBuffer(jsonData))

	request.Header = defaultHeader()

	if DEBUG {
		dump, _ := httputil.DumpRequest(request, true)

		fmt.Printf("auth request: %q\n", dump)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if DEBUG {
		fmt.Printf("auth response: %s\n", body)
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data["session_id"], nil
}

func createAuthHeader(s string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s))
}
