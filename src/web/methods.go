package web

import (
	"bufio"
	"configs"
	"encoding/json"
	"fmt"
	"github.com/ceph/go-ceph/rados"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"metadata"
	"net/http"
	"strconv"
	"strings"
	"wrados"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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

		ss, eror := metadata.RedClient(pool+"/"+name, "get", "")
		filez := []string{}

		_, infostat := r.URL.Query()["info"]
		if infostat {
			if lo != nil {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(lo.Error()))
				_, _ = w.Write([]byte("\n"))
			} else {
				if eror == nil {
					ns := strings.Split(ss, ",")
					numbSegments := strconv.Itoa(len(ns))

					fileInfo := &FileInfo{
						Size:     ns[len(ns)-1],
						Pool:     pool,
						Segments: numbSegments,
						Name:     name,
					}
					b, _ := json.Marshal(fileInfo)
					_, _ = w.Write(b)
					_, _ = w.Write([]byte("\n"))
				} else {
					fileInfo := &FileInfo{
						Size:     strconv.Itoa(int(xo.Size)),
						Pool:     pool,
						Segments: "1",
						Name:     name,
					}
					b, _ := json.Marshal(fileInfo)
					_, _ = w.Write(b)
					_, _ = w.Write([]byte("\n"))

				}
			}

		} else {
			readFilez := func(name string) {
				if lo == nil {
					of := uint64(0)
					mx := uint64(20480)
					if xo.Size-of < mx {
						mx = xo.Size
					}

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
					wrados.Writelog(r.Method, xo.Size, "bytes", r.URL, "from", pool)
				} else {
					_, _ = w.Write([]byte(respCodewriter(lo, w, r)))
				}
			}

			if eror != nil {
				filez = append(filez, name)
				for file := range filez {
					w.Header().Set("Content-Length", strconv.FormatUint(xo.Size, 10))
					readFilez(filez[file])
				}
			} else {
				var fsize uint64
				fileparts := strings.Split(ss, ",")
				fileparts = fileparts[:len(fileparts)-1]
				for filepart := range fileparts {
					name = fileparts[filepart]
					xo, _ = ioctx.Stat(name)
					fsize = fsize + xo.Size
				}
				w.Header().Set("Content-Length", strconv.FormatUint(fsize, 10))
				for filepart := range fileparts {
					name = fileparts[filepart]
					xo, _ = ioctx.Stat(name)
					readFilez(fileparts[filepart])
				}
			}
		}

	} else {
		wrados.Writelog("Pool " + pool + " does not exists")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 Internal Server Error \n"))
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
					xxx := 0
					size := 0
					fileSegments := make([]string, 0)
					for {

						line, err := reqBody.ReadBytes('\n')
						mukuch = append(mukuch, line...)
						if err == io.EOF {
							name := name + "-" + randSeq(10)
							xo, _ := ioct.Stat(name)
							_ = ioct.Write(name, mukuch, xo.Size)

							fileSegments = append(fileSegments, name)
							xxx = xxx + 1

							break
						}
						if len(mukuch) > configs.Conf.Uploadmaxpart {
							name := name + "-" + randSeq(10)
							xo, _ := ioct.Stat(name)
							_ = ioct.Write(name, mukuch, xo.Size)

							fileSegments = append(fileSegments, name)
							xxx = xxx + 1
							lenMukuch := len(mukuch)

							size = size + lenMukuch
							wrados.Writelog(r.Method, lenMukuch, "bytes, segment", name, "of", r.URL, "to", pool)
							mukuch = nil
						}

						//fileSegments = append(fileSegments, string(size))
					}
					fileSegments = append(fileSegments, strconv.Itoa(size))
					log.Println("Created File", name, "In", pool)
					fmeta := strings.Join(fileSegments, ",")
					_, err := metadata.RedClient(pool+"/"+name, "set", fmeta)
					if err != nil {
						log.Println("error setting metadata:", err)
					}

				}

				wrados.Writelog(r.Method, r.Header.Get("Content-Length"), "bytes", r.URL, "to", pool)
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
		msg := "Server is running in read only mode ! \n"
		wrados.Writelog(msg)
		_, _ = w.Write([]byte(msg))
	}
}

func Del(w http.ResponseWriter, r *http.Request) {
	switch configs.Conf.DangeZone {
	case true:

		s := strings.Split(r.URL.Path, "/")
		if len(s) >= 3 {
			pool := s[1]
			name := strings.Join(s[2:], "/")

			ss, eror := metadata.RedClient(pool+"/"+name, "get", "")
			filez := []string{}

			if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
				randindex := rand.Intn(len(wrados.Rconnect.Connection))
				ioct, _ := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)

				filez = append(filez, name)
				if eror == nil {
					//filez = append(filez, name)

					fileparts := strings.Split(ss, ",")
					fileparts = fileparts[:len(fileparts)-1]
					for filepart := range fileparts {
						filez = append(filez, fileparts[filepart])
					}

				}

				for filename := range filez {
					f := ioct.Delete(filez[filename])
					if f != nil {
						_, _ = fmt.Fprintf(w, respCodewriter(f, w, r))
					} else {
						wrados.Writelog(r.Method, filez[filename], "from", pool)
						msg := http.StatusText(200) + ", Deleted: " + r.URL.String() + "\n"
						_, _ = w.Write([]byte(msg))
					}
				}
				_, _ = metadata.RedClient(pool+"/"+name, "del", "")

			}
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		msg := "Dangerous commands are disabled ! \n"
		wrados.Writelog(msg)
		_, _ = w.Write([]byte(msg))
	}
}
