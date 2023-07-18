package session

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Response is a wrapper to a http.Response pointer, with the isolated Cookies and response data as a []byte
type Response struct {
	*http.Response
	// DataBytes response data as []byte
	DataBytes []byte
	// Cookies response cookies
	Cookies []*http.Cookie
}

// Session is the TracimAPI client
type Session struct {
	// config is the configuration for the api
	config
	isAuth bool
	UserID string
	// Session contains the session cookie
	Session      http.Cookie
	eventChannel chan TracimLiveMessage
	eventHandler map[string]TLMHandler
	errorChannel chan error
}

// GenerateRequest Generates an authenticated Request on the Tracim api, if there is no data to send, set it to nil
func (s *Session) GenerateRequest(method, endpoint string, data []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", s.Url, endpoint)
	if data == nil {
		data = []byte("")
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if s.isAuth {
		req.AddCookie(&s.Session)
	}

	return req, err
}

// Request Does an authenticated Request on the Tracim api, if there is no data to send, set it to nil
func (s *Session) Request(method, endpoint string, data []byte) (Response, error) {
	if s.isAuth && !time.Now().Before(s.Session.Expires) {
		err := s.Auth()
		if err != nil {
			return Response{}, fmt.Errorf("auth - %s", err.Error())
		}
	}
	client := &http.Client{}

	req, err := s.GenerateRequest(method, endpoint, data)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Response{}, fmt.Errorf(resp.Status)
	}
	if err != nil {
		return Response{}, err
	}
	return Response{resp, body, resp.Cookies()}, nil
}

// SetCredentials set the credentials for the TracimAPI client
func (s *Session) SetCredentials(credentials Credentials) *Session {
	s.Credentials = credentials
	return s
}

// New Creates a TracimAPI client, takes the link to the Tracim api, usually https://example.com/api
func New(link string) (s *Session) {
	s = &Session{
		config:       config{Url: link},
		eventHandler: make(map[string]TLMHandler),
		eventChannel: make(chan TracimLiveMessage),
		errorChannel: make(chan error),
	}

	return s
}
