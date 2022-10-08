package TracimAPI

import (
	"bytes"
	"fmt"
	"github.com/Millefeuille42/TracimAPI/Config"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Data    []byte
	Cookies []*http.Cookie
}

type Api struct {
	Config.Config
	isAuth  bool
	Session http.Cookie
}

func (a *Api) Request(method, endpoint string, data []byte) (Response, error) {
	if a.isAuth && !time.Now().Before(a.Session.Expires) {
		err := a.Auth()
		if err != nil {
			return Response{}, fmt.Errorf("auth - %s", err.Error())
		}
	}
	url := fmt.Sprintf("%s%s", a.Url, endpoint)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return Response{}, err
	}
	if a.isAuth {
		req.AddCookie(&a.Session)
	}
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
	resp.Cookies()
	return Response{body, resp.Cookies()}, nil
}

func (a *Api) SetCredentials(mail, password string) *Api {
	a.Mail = mail
	a.Password = password
	return a
}

func New(link string) *Api {
	return &Api{Config: Config.Config{Url: link}}
}
