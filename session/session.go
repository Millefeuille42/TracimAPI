package session

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	*http.Response
	DataBytes []byte
	Cookies   []*http.Cookie
}

type Session struct {
	config
	isAuth       bool
	Session      http.Cookie
	eventChannel chan TracimLiveMessage
	eventHandler map[string]TLMHandler
	errorChannel chan error
}

func (a *Session) GenerateRequest(method, endpoint string, data []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", a.Url, endpoint)
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if a.isAuth {
		req.AddCookie(&a.Session)
	}

	return req, err
}

func (a *Session) Request(method, endpoint string, data []byte) (Response, error) {
	if a.isAuth && !time.Now().Before(a.Session.Expires) {
		err := a.Auth()
		if err != nil {
			return Response{}, fmt.Errorf("auth - %s", err.Error())
		}
	}
	client := &http.Client{}

	req, err := a.GenerateRequest(method, endpoint, data)
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

func (a *Session) SetCredentials(credentials Credentials) *Session {
	a.Credentials = credentials
	return a
}

func New(link string) (s *Session) {
	s = &Session{
		config:       config{Url: link},
		eventHandler: make(map[string]TLMHandler),
		eventChannel: make(chan TracimLiveMessage),
		errorChannel: make(chan error),
	}

	return s
}
