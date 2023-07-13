package session

type Credentials struct {
	Mail     string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type config struct {
	Credentials
	Url string
}
