package configs

import (
	"net/http"
)

func GetIP(r *http.Request) string {
	switch r.Header.Get("X-Forwarded-For") {
	case "":
		return r.RemoteAddr
	default:
		return r.Header.Get("X-Forwarded-For")

	}
}
