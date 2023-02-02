package web

import (
	"bufio"
	"configs"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"wrados"
)

var users = map[string]string{
	//"test":  "secret",
}

func isAuthorised(username, password string) bool {

	md5HashInBytes := md5.Sum([]byte(password))
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])

	pass, ok := users[username]

	if !ok {
		return false
	}
	return md5HashInString == pass
}

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	//w.Header().Add("Content-Type", "application/json")
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Authentication Required"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("No basic auth present\n"))
		wrados.Writelog("401 Unauthorized: No basic auth present")
		return false
	}

	if !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid Credentials\n"))
		wrados.Writelog("401 Unauthorized: Invalid Credentials")
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true

}

func PopulateUsers() {
	for {
		c, e := os.Open("users.txt")
		content := bufio.NewScanner(c)
		if e != nil {
			wrados.Writelog(e)
		}

		for content.Scan() {
			if len(content.Text()) > 0 {
				z := strings.Split(content.Text(), " ")
				if _, ok := users[z[0]]; !ok {
					wrados.Writelog("Found new user: " + z[0] + ", enabling!")
					users[z[0]] = z[1]
				}
			}
			//fmt.Println(content.Text())
		}
		_ = c.Close()
		time.Sleep(10 * time.Second)
	}
}

func dynHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.Header, r.Method)
	switch r.Method {
	case "GET":
		if configs.Conf.AuthRead {
			if authenticate(w, r) {
				momo.incrementGet()
				Get(w, r)
			}
		} else {
			momo.incrementGet()
			Get(w, r)
		}
	case "POST", "PUT":
		_, ko := r.Header["Content-Length"]
		if !ko {
			wrados.Writelog("Header \"Content-Length\" is not present in request, aborting")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("400: Header \"Content-Length\" is not present in request, aborting \n"))
			return
		}
		if configs.Conf.AuthWrite {
			if authenticate(w, r) {
				momo.incrementPost()
				Put(w, r)
			}
		} else {
			momo.incrementPost()
			Put(w, r)
		}
	case "DELETE":
		if configs.Conf.AuthWrite {
			if authenticate(w, r) {
				momo.incrementDel()
				Del(w, r)
			}
		} else {
			momo.incrementDel()
			Del(w, r)
		}
	//	momo.incrementHead()
	//	Head(w, r)
	default:
		_, _ = fmt.Fprintf(w, "Sorry, only GET, POST, PUT, HEAD and DELETE methods are supported.\n")
	}
}

func playmux0() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dynHandler)

	s1 := &http.Server{
		Addr:         configs.Conf.HttpAddress,
		Handler:      mux,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 10800 * time.Second,
		IdleTimeout:  1,
	}
	_ = s1.ListenAndServe()
}

func mxhandl(w http.ResponseWriter, _ *http.Request) {
	mz := printStats()
	_, _ = fmt.Fprintln(w, mz)
}

func playmux1() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", mxhandl)
	//users[configs.Conf.ServerUser] = configs.Conf.ServerPass

	s2 := http.Server{
		Addr:         configs.Conf.MonAddress,
		Handler:      mux1,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	wrados.Writelog("Starting monitoring instance at:", configs.Conf.MonAddress)
	_ = s2.ListenAndServe()

}

func RunServer() {
	if configs.Conf.Monenabled {
		go playmux1()
	}
	wrados.Writelog("Starting WebRados server at:", configs.Conf.HttpAddress)
	runtime.Gosched()
	playmux0()
}
