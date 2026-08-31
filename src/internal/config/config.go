package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DbAuth        string
	JWTkey        []byte
	JWTaccessTTL  time.Duration
	JWTrefreshTTL time.Duration
	ServerAddr    string
}

func ParseConfig() (*Config, error) {
	var err error
	getStr := func(key string) string {
		if err != nil {
			return ""
		}
		str := os.Getenv(key)
		if str == "" {
			err = fmt.Errorf("Env var not found: %q", key)
		}
		return str
	}
	getDur := func(key string) time.Duration {
		str := getStr(key)
		if err != nil {
			return time.Second
		}
		var dur time.Duration
		dur, err = time.ParseDuration(str)
		return dur
	}

	conf := Config{
		DbAuth:        getStr("DB_AUTH"),
		JWTkey:        []byte(getStr("JWT_KEY")),
		ServerAddr:    getStr("SERVER_ADDR") + ":" + getStr("SERVER_PORT"),
		JWTaccessTTL:  getDur("JWT_ACCESS_TTL"),
		JWTrefreshTTL: getDur("JWT_REFRESH_TTL"),
	}

	return &conf, err
}
