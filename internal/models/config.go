package models

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn string
	}
	Limiter struct {
		Enabled bool
		Rps     float64
		Burst   int
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}
