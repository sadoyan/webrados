package configs

import (
	"net/http"
)

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetIP(r *http.Request) string {
	switch r.Header.Get("X-Forwarded-For") {
	case "":
		return r.RemoteAddr
	default:
		return r.Header.Get("X-Forwarded-For")

	}
}
