package redlock

type Config struct {
	ExpirySeconds int
	Redis         struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}
