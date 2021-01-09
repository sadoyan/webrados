package web

import (
	"bufio"
	"configs"
	"fmt"
	"github.com/ceph/go-ceph/rados"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wrados"
)

func respCodewriter(f error, w http.ResponseWriter, r *http.Request) string {
	if strings.Split(f.Error(), ",")[1] == " No such file or directory" {
		w.WriteHeader(http.StatusNotFound)
		log.Println(r.Method, f.Error(), r.URL.String())
		return http.StatusText(404) + ": " + r.URL.String() + "\n"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(r.Method, f.Error(), r.URL.String())
		return http.StatusText(500) + ": " + r.URL.String() + "\n"
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]
	name := strings.Join(s[2:], "/")
	if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
		ioctx, e := wrados.Rconnect.Connection.OpenIOContext(pool)
		if e != nil {
			log.Println(e)
		}
		xo, lo := ioctx.Stat(name)
		if lo == nil {
			of := uint64(0)
			mx := uint64(4096)
			if xo.Size-of < mx {
				mx = xo.Size
			}
			w.Header().Set("Content-Length", strconv.FormatUint(xo.Size, 10))
			for {
				if xo.Size-of <= mx {
					mx = xo.Size - of
				}
				bytesOut := make([]byte, mx)
				_, err := ioctx.Read(name, bytesOut, of)
				if err != nil {
					log.Println(err)
				}
				_, er := w.Write(bytesOut)
				if er != nil {
					log.Println(er)
				}
				of = of + mx
				if of >= xo.Size {
					break
				}
			}
			log.Println("Method", r.Method, xo.Size, "bytes", name, "from", pool)
		} else {
			_, _ = w.Write([]byte(respCodewriter(lo, w, r)))
		}
	} else {
		log.Println("Pool " + pool + " does not exists")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 Internal Server Error \n"))
	}
}

func Put(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	if len(s) >= 3 {
		if _, ok := wrados.Rconnect.Poolnames[s[1]]; ok {
			pool := s[1]
			name := strings.Join(s[2:], "/")

			ioct, _ := wrados.Rconnect.Connection.OpenIOContext(pool)
			lenq, _ := strconv.Atoi(r.Header.Get("Content-Length"))

			if lenq < configs.Conf.Uploadmaxpart {
				reqBody, _ := ioutil.ReadAll(r.Body)
				_ = ioct.Write(name, reqBody, 0)
			} else {
				reqBody := bufio.NewReader(r.Body)
				_ = ioct.Create(name, rados.CreateOption(lenq))
				mukuch := make([]byte, 0)
				for {
					line, err := reqBody.ReadBytes('\n')
					mukuch = append(mukuch, line...)
					if err == io.EOF {
						xo, _ := ioct.Stat(name)
						_ = ioct.Write(name, mukuch, xo.Size)
						break
					}
					if len(mukuch) > configs.Conf.Uploadmaxpart {
						xo, _ := ioct.Stat(name)
						_ = ioct.Write(name, mukuch, xo.Size)
						mukuch = nil
					}
				}
			}
			log.Println("Method", r.Method, r.Header.Get("Content-Length"), "bytes", name, "to", pool)
		} else {
			log.Println("Invalid pool name")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500: Invalid pool name \n"))
		}

	} else {
		log.Println("File path is too short")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500: File path is too short \n"))
	}
}

func Del(w http.ResponseWriter, r *http.Request) {
	if configs.Conf.DangeZone {
		s := strings.Split(r.URL.Path, "/")
		if len(s) >= 3 {
			if _, ok := wrados.Rconnect.Poolnames[s[1]]; ok {
				pool := s[1]
				name := strings.Join(s[2:], "/")
				ioct, _ := wrados.Rconnect.Connection.OpenIOContext(pool)
				f := ioct.Delete(name)
				if f != nil {
					_, _ = fmt.Fprintf(w, respCodewriter(f, w, r))
				} else {
					log.Println("Method", r.Method, name, "from", pool)
					msg := http.StatusText(200) + ", Deleted: " + r.URL.String() + "\n"
					_, _ = fmt.Fprintf(w, msg)
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		msg := "Dangerous commands are disabled !"
		log.Println(msg)
		_, _ = fmt.Fprintf(w, msg+"\n")
	}
}

func Head(w http.ResponseWriter, r *http.Request) {
	if configs.Conf.DangeZone {
		s := strings.Split(r.URL.Path, "/")
		if len(s) == 2 {
			if _, ok := wrados.Rconnect.Poolnames[s[1]]; ok {
				pool := s[1]
				m, _ := wrados.Rconnect.Connection.OpenIOContext(pool)
				c, _ := m.GetPoolStats()
				fmt.Println(c)
			}
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		msg := "Dangerous commands are disabled !"
		log.Println(msg)
		_, _ = fmt.Fprintf(w, msg+"\n")
	}
}
