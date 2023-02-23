package web

import (
	"auth"
	"configs"
	"fmt"
	"net/http"
	"runtime"
	"time"
	"wrados"
)

func dynHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if configs.Conf.AuthRead {
			if auth.CheckAuth(w, r) {
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
			if auth.CheckAuth(w, r) {
				momo.incrementPost()
				Put(w, r)
			}
		} else {
			momo.incrementPost()
			Put(w, r)
		}
	case "DELETE":
		if configs.Conf.AuthWrite {
			if auth.CheckAuth(w, r) {
				momo.incrementDel()
				Del(w, r)
			}
		} else {
			momo.incrementDel()
			Del(w, r)
		}
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
