package Config

type Credentials struct {
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type Config struct {
	Credentials
	Url string
}
