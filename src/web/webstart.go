package web

import (
	"configs"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

var users = map[string]string{
	//"test":  "secret",
}

func isAuthorised(username, password string) bool {
	pass, ok := users[username]
	if !ok {
		return false
	}
	return password == pass
}

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	//w.Header().Add("Content-Type", "application/json")
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Authentication Required"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("No basic auth present\n"))
		log.Println("401 Unauthorized: No basic auth present")
		return false
	}

	if !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid Credentials\n"))
		log.Println("401 Unauthorized: Invalid Credentials")
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true

}

// -------------------------------------------------------------------------- //
func dynHandler(w http.ResponseWriter, r *http.Request) {
	users[configs.Conf.ServerUser] = configs.Conf.ServerPass
	//fmt.Println(users)
	switch r.Method {
	case "GET":
		if configs.Conf.AuthRead {
			if authenticate(w, r) {
				Get(w, r)
			}
		} else {
			Get(w, r)
		}
	case "POST", "PUT":
		if configs.Conf.AuthWrite {
			if authenticate(w, r) {
				Put(w, r)
			}
		} else {
			Put(w, r)
		}

	case "DELETE":
		if configs.Conf.AuthWrite {
			if authenticate(w, r) {
				Del(w, r)
			}
		} else {
			Del(w, r)
		}
		//Del(w, r)
	case "HEAD":
		Head(w, r)
	default:
		_, _ = fmt.Fprintf(w, "Sorry, only GET, HEAD, POST, PUT and DELETE methods are supported.\n")
	}
}

func playmux0() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dynHandler)

	s1 := http.Server{
		Addr:         configs.Conf.HttpAddress,
		Handler:      mux,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	_ = s1.ListenAndServe()
}

func mxhandl(w http.ResponseWriter, r *http.Request) {
	mz := printStats()
	_, _ = fmt.Fprintln(w, mz)
}

func playmux1() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", mxhandl)
	s2 := http.Server{
		Addr:         configs.Conf.MonAddress,
		Handler:      mux1,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	log.Println("Starting monitoring instance at:", configs.Conf.MonAddress)
	_ = s2.ListenAndServe()

}

func RunServer() {
	if configs.Conf.Monenabled {
		go playmux1()
	}
	log.Println("Starting WebRados server at:", configs.Conf.HttpAddress)
	runtime.Gosched()
	playmux0()
}
