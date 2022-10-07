package Config

type Credentials struct {
	Mail     string
	Password string
}

type Config struct {
	Credentials
	Url string
}
