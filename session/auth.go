package session

import "encoding/json"

type T struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Session) Auth() error {
	a.isAuth = false
	data, err := json.Marshal(a.Credentials)
	if err != nil {
		return err
	}
	response, err := a.Request("POST", "/auth/login", data)
	if err != nil {
		return err
	}
	for _, cookie := range response.Cookies {
		if cookie.Name == "session_key" {
			a.Session = *cookie
			break
		}
	}
	a.isAuth = true
	return nil
}
