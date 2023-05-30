package web

import (
	"auth"
	"configs"
	"fmt"
	"io"
	"metadata"
	"net/http"
	"runtime"
	"time"
	"tools"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	/*
		curl -s  -XPOST -H "X-API-KEY: $B" --data-binary @/tmp/oo.json  'http://192.168.111.2:8080/admin?sign' | python3 -mjson.tool

		REQUEST
		{
		    "http://192.168.111.2:8080/bublics/a23fa6e7-7b1d-4172-a48d-7e143d798788.jpeg": 180,
		    "http://192.168.111.2:8080/bublics/c2cccb2f-64c7-4fed-8396-52246b962b79.jpeg": 120,
		    "http://192.168.111.2:8080/bublics/dca6056b-49f9-478a-bfa8-bf61a4b2d89a.jpeg": 90,
		    "http://192.168.111.2:8080/bublics/3379d4ae-9648-4492-815b-8dcb6bd2bc13.jpeg": 150
		}
		RESPONSE
		{
		    "http://192.168.111.2:8080/bublics/3379d4ae-9648-4492-815b-8dcb6bd2bc13.jpeg": "http://192.168.111.2:8080/bublics/3379d4ae-9648-4492-815b-8dcb6bd2bc13.jpeg?expiry=1685589884&signature=5Hl0oosDIWV9TcuBNlYqICu5-kCGCp5nwIWEdsipx60",
		    "http://192.168.111.2:8080/bublics/a23fa6e7-7b1d-4172-a48d-7e143d798788.jpeg": "http://192.168.111.2:8080/bublics/a23fa6e7-7b1d-4172-a48d-7e143d798788.jpeg?expiry=1685440064&signature=0vGnOQhYmUdLZTfh3zIUcJ5dgsiCxINI_5ROOW9K5ew",
		    "http://192.168.111.2:8080/bublics/c2cccb2f-64c7-4fed-8396-52246b962b79.jpeg": "http://192.168.111.2:8080/bublics/c2cccb2f-64c7-4fed-8396-52246b962b79.jpeg?expiry=1685440004&signature=0dK-LTFhdQ_LNPG75GOgrDfeldWsoT-VjBqx4bS6a3o",
		    "http://192.168.111.2:8080/bublics/dca6056b-49f9-478a-bfa8-bf61a4b2d89a.jpeg": "http://192.168.111.2:8080/bublics/dca6056b-49f9-478a-bfa8-bf61a4b2d89a.jpeg?expiry=1685439974&signature=OcoNRY7q3H3W_M1eZrfS99DHfHeJTkR8Q0eTcUdTjgs"
		}
	*/

	momo.incrementGet()
	if !auth.DoAdminAuth(r) {
		http.Error(w, http.StatusText(401), 401)
		tools.WriteLogs(tools.GetIP(r), r.Method, "401 Unauthorized", r.URL)
		return
	}
	_, getjwt := r.URL.Query()["genjwt"]
	if getjwt {
		k, _ := io.ReadAll(r.Body)
		tok, _ := auth.GenJWTtoken(k)
		_, _ = w.Write(tok)
		_, _ = w.Write([]byte("\n"))
		tools.WriteLogs("Successfully generated JWT token:", tools.GetIP(r), r.URL)
		return
	}
	_, urlsig := r.URL.Query()["sign"]
	if urlsig {
		ret := auth.SignUrl(r.Body)
		_, _ = w.Write(ret)
		_, _ = w.Write([]byte("\n"))
		tools.WriteLogs("Successfully Signed bunch of urls:", tools.GetIP(r), r.URL)
		return
	}
	_, purgecache := r.URL.Query()["purgecache"]
	if purgecache {
		_ = metadata.Cache.Reset()
		_ = metadata.Cache.ResetStats()
		tools.WriteLogs(tools.GetIP(r), r.Method, "Purging everything from cache")
		return
	}
	_, cachestats := r.URL.Query()["purgecachestats"]
	if cachestats {
		_ = metadata.Cache.ResetStats()
		tools.WriteLogs(tools.GetIP(r), r.Method, "Resetting cache statistics")
		return
	}

}

func dynHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if configs.Conf.AuthRead {
			if auth.DoAuth(r) {
				momo.incrementGet()
				Get(w, r)
			} else {
				momo.incrementGet()
				http.Error(w, http.StatusText(401), 401)
				tools.WriteLogs(tools.GetIP(r), r.Method, "401 Unauthorized", r.URL)
			}
		} else {
			momo.incrementGet()
			Get(w, r)
		}
	case "POST", "PUT":
		/*
			_, getjwt := r.URL.Query()["genjwt"]
			if getjwt {
				// curl -XPOST -H "X-API-KEY: $B" -d '{"username": "valog","password": "guggush","exp": 1685110930}'  http://192.168.111.2:8080/?genjwt
				if auth.DoAdminAuth(r) {
					momo.incrementGet()
					k, _ := io.ReadAll(r.Body)
					tok, _ := auth.GenJWTtoken(k)
					_, _ = w.Write(tok)
					_, _ = w.Write([]byte("\n"))
					tools.WriteLogs("Sucesfully generated JWT token:", tools.GetIP(r), r.URL)
					return
				} else {
					momo.incrementGet()
					http.Error(w, http.StatusText(401), 401)
					tools.WriteLogs(tools.GetIP(r), r.Method, "401 Unauthorized", r.URL)
					return
				}
			}
		*/

		_, ko := r.Header["Content-Length"]
		if !ko {
			tools.WriteLogs("Header \"Content-Length\" is not present in request, aborting")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("400: Header \"Content-Length\" is not present in request, aborting \n"))
			return
		}
		if configs.Conf.AuthWrite {
			if auth.DoAuth(r) {
				momo.incrementPost()
				Put(w, r)
			} else {
				momo.incrementGet()
				http.Error(w, http.StatusText(401), 401)
				tools.WriteLogs(tools.GetIP(r), r.Method, "401 Unauthorized", r.URL)
			}
		} else {
			momo.incrementPost()
			Put(w, r)
		}
	case "DELETE":
		if configs.Conf.AuthWrite {
			if auth.DoAuth(r) {
				momo.incrementDel()
				Del(w, r)
			} else {
				momo.incrementDel()
				http.Error(w, http.StatusText(401), 401)
				tools.WriteLogs(tools.GetIP(r), r.Method, "401 Unauthorized", r.URL)
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
	mux.HandleFunc("/.admin", adminHandler)

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
	tools.WriteLogs("Starting monitoring instance at:", configs.Conf.MonAddress)
	_ = s2.ListenAndServe()

}

func RunServer() {
	if configs.Conf.Monenabled {
		go playmux1()
	}
	tools.WriteLogs("Starting WebRados server at:", configs.Conf.HttpAddress)
	runtime.Gosched()
	playmux0()
}
