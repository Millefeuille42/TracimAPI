package session

import "encoding/json"

type authResponse struct {
	UserID string `json:"user_Id"`
}

// Auth authenticate client to the Tracim api
func (s *Session) Auth() error {
	s.isAuth = false
	data, err := json.Marshal(s.Credentials)
	if err != nil {
		return err
	}
	response, err := s.Request("POST", "/auth/login", data)
	if err != nil {
		return err
	}
	for _, cookie := range response.Cookies {
		if cookie.Name == "session_key" {
			s.Session = *cookie
			break
		}
	}

	respData := authResponse{}
	err = json.Unmarshal(response.DataBytes, &respData)
	if err != nil {
		return err
	}

	s.userID = respData.userID
	s.isAuth = true
	return nil
}
