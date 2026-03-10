package util

import (
	"os"
)

func Port() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8080"
	}
	return p
}

func BaseURL() string {
	return "http://localhost:" + Port()
}

func ApiURL(path string) string {
	return BaseURL() + path
}
