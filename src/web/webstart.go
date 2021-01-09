package web

import (
	"configs"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

//func Startserver() {
//	fmt.Println("Hi")
//}
// -------------------------------------------------------------------------- //
func dynHandler(w http.ResponseWriter, r *http.Request) {
	//const unauth = http.StatusUnauthorized
	//if configs.Conf.serverAuth {
	//	auth := r.Header.Get("Authorization")
	//	if !strings.HasPrefix(auth, "Basic ") {
	//		log.Print("Invalid authorization:", auth)
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		return
	//	}
	//	up, err := base64.StdEncoding.DecodeString(auth[6:])
	//	if err != nil {
	//		log.Print("authorization decode error:", err)
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		return
	//	}
	//	if string(up) != authorized["server"] {
	//		http.Error(w, http.StatusText(unauth), unauth)
	//		return
	//	}
	//}
	// -- ---------- -- //

	switch r.Method {
	case "GET":
		Get(w, r)
	case "POST", "PUT":
		Put(w, r)
	case "DELETE":
		Del(w, r)
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
	_ = s2.ListenAndServe()

}

func RunServer() {
	//http.HandleFunc("/", dynHandler)
	//fmt.Println("starting server at: " + configs.Conf.HttpAddress)

	if configs.Conf.Monenabled {
		go playmux1()
	}

	log.Print("Started WebRados ")
	runtime.Gosched()

	playmux0()

	//if configs.Conf.Monenabled {
	//	go playmux1()
	//}

	//log.Println("Currently running", runtime.NumGoroutine(), "Goroutines")

	//forever := make(chan bool)
	//<-forever

	//log.Fatal(http.ListenAndServe(":9090", nil))

}
