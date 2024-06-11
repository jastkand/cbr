package config

import (
	"os"
)

var Port string

func Load() {
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "5000"
	}
}
