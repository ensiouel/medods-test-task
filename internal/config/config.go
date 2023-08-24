package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Server  Server
	Mongo   Mongo
	Session Session
}

type Server struct {
	Addr string `env:"SERVER_ADDR"`
}

type Mongo struct {
	Host string `env:"MONGO_HOST"`
	Port string `env:"MONGO_PORT"`
}

type Session struct {
	AccessTokenExpiration  time.Duration `env:"SESSION_ACCESS_TOKEN_EXPIRATION"`
	RefreshTokenExpiration time.Duration `env:"SESSION_REFRESH_TOKEN_EXPIRATION"`
	Secret                 []byte        `env:"SESSION_SECRET"`
}

func New() Config {
	var config Config
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
