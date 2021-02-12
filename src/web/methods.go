package web

import (
	"bufio"
	"configs"
	"fmt"
	"github.com/ceph/go-ceph/rados"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"wrados"
)

func respCodewriter(f error, w http.ResponseWriter, r *http.Request) string {
	if strings.Split(f.Error(), ",")[1] == " No such file or directory" {
		w.WriteHeader(http.StatusNotFound)
		wrados.Writelog(r.Method, f.Error(), r.URL.String())
		return http.StatusText(404) + ": " + r.URL.String() + "\n"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		wrados.Writelog(r.Method, f.Error(), r.URL.String())
		return http.StatusText(500) + ": " + r.URL.String() + "\n"
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]
	name := strings.Join(s[2:], "/")
	if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
		randindex := rand.Intn(len(wrados.Rconnect.Connection))
		ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
		if e != nil {
			wrados.Writelog(e)
		}
		xo, lo := ioctx.Stat(name)
		if lo == nil {
			of := uint64(0)
			mx := uint64(20480)
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
					wrados.Writelog(err)
					break
				}
				_, er := w.Write(bytesOut)
				if er != nil {
					wrados.Writelog(er)
					break
				}
				of = of + mx
				if of >= xo.Size {
					break
				}
			}
			wrados.Writelog("Method", r.Method, xo.Size, "bytes", name, "from", pool)
		} else {
			_, _ = w.Write([]byte(respCodewriter(lo, w, r)))
		}
	} else {
		wrados.Writelog("Pool " + pool + " does not exists")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 Internal Server Error \n"))
	}
}

func Got(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]
	name := strings.Join(s[2:], "/")
	if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
		randindex := rand.Intn(len(wrados.Rconnect.Connection))
		ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
		if e != nil {
			fmt.Println(e)
		}
		xo, _ := ioctx.Stat(name)
		bytesOut := make([]byte, xo.Size)
		out, _ := ioctx.Read(name, bytesOut, 0)
		fmt.Println(out, pool, name, xo.Size)
		_, _ = w.Write(bytesOut)
	} else {
		fmt.Println("Pool " + pool + " does not exists")
	}
}

func Put(w http.ResponseWriter, r *http.Request) {
	switch configs.Conf.Readonly {
	case false:
		s := strings.Split(r.URL.Path, "/")
		if len(s) >= 3 {
			pool := s[1]
			name := strings.Join(s[2:], "/")
			if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
				randindex := rand.Intn(len(wrados.Rconnect.Connection))
				ioct, _ := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
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
				wrados.Writelog("Method", r.Method, r.Header.Get("Content-Length"), "bytes", name, "to", pool)
			} else {
				wrados.Writelog("Invalid pool name")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("500: Invalid pool name \n"))
			}

		} else {
			wrados.Writelog("File path is too short")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500: File path is too short \n"))
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		msg := "Server is running in read only mode !"
		wrados.Writelog(msg)
		_, _ = fmt.Fprintf(w, msg+"\n")

	}

}

func Del(w http.ResponseWriter, r *http.Request) {
	switch configs.Conf.DangeZone {
	case true:
		s := strings.Split(r.URL.Path, "/")
		if len(s) >= 3 {
			pool := s[1]
			name := strings.Join(s[2:], "/")
			if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
				randindex := rand.Intn(len(wrados.Rconnect.Connection))
				ioct, _ := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
				f := ioct.Delete(name)
				if f != nil {
					_, _ = fmt.Fprintf(w, respCodewriter(f, w, r))
				} else {
					wrados.Writelog("Method", r.Method, name, "from", pool)
					msg := http.StatusText(200) + ", Deleted: " + r.URL.String() + "\n"
					_, _ = fmt.Fprintf(w, msg)
				}
			}
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		msg := "Dangerous commands are disabled !"
		wrados.Writelog(msg)
		_, _ = fmt.Fprintf(w, msg+"\n")
	}
}

func Head(w http.ResponseWriter, r *http.Request) {
	if configs.Conf.DangeZone {
		s := strings.Split(r.URL.Path, "/")
		if len(s) == 2 {
			pool := s[1]
			if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
				randindex := rand.Intn(len(wrados.Rconnect.Connection))
				m, _ := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
				c, _ := m.GetPoolStats()
				fmt.Println(c)
			}
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		msg := "Dangerous commands are disabled !"
		wrados.Writelog(msg)
		_, _ = fmt.Fprintf(w, msg+"\n")
	}
}
