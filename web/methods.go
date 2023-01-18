package web

import (
	"bufio"
	"configs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"metadata"
	"net/http"
	"strconv"
	"strings"
	"wrados"

	"github.com/ceph/go-ceph/rados"
)

//var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
//func randSeq(n int) string {
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(b)
//}

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

func Split(r rune) bool {
	return r == '=' || r == '-'
}

func readFile(w http.ResponseWriter, r *http.Request, name string, pool string, xo rados.ObjectStat, of uint64) bool {
	randindex := rand.Intn(len(wrados.Rconnect.Connection))
	ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)

	if e != nil {
		wrados.Writelog(e)
	}
	mx := uint64(102400)
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
			if !strings.HasPrefix(er.Error(), "write tcp") {
				wrados.Writelog(er)
			}
			return false

		}
		of = of + mx
		if of >= xo.Size {
			break
		}
	}
	wrados.Writelog(r.Method, xo.Size, "bytes", name, "from", pool)
	return true
}

func Get(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]

	var minrange int
	var contentlenght int
	var of uint64
	switch pool {
	case "favicon.ico":
		// DO nothing !
	default:
		name := strings.Join(s[len(s)-1:], "/")
		extension := strings.Split(name, ".")[1]

		if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
			randindex := rand.Intn(len(wrados.Rconnect.Connection))
			ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
			if e != nil {
				wrados.Writelog(e)
			}
			filename, eror := metadata.DBClient(pool+"/"+name, "get", "")
			xo, lo := ioctx.Stat(name)
			if lo != nil {
				errormsg := strings.Split(fmt.Sprint(lo), ",")
				wrados.Writelog(errormsg)
				switch errormsg[len(errormsg)-1] {
				case " No such file or directory":
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte("404 File not found \n"))
				default:
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("500 Internal Server Error \n"))
				}
				return
			}

			_, infostat := r.URL.Query()["info"]

			if infostat {
				if lo != nil {
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte(lo.Error()))
					_, _ = w.Write([]byte("\n"))
				} else {
					if eror == nil {
						ns := strings.Split(filename, ",")
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
				return
			}

			mime, mok := HttpMimes.Videos[extension]
			switch mok {
			case false:
				w.Header().Set("Content-Length", strconv.FormatUint(uint64(xo.Size), 10))
				readFile(w, r, name, pool, xo, of)
			case true:
				var fsize uint64
				var fileparts []string
				xo, _ = ioctx.Stat(name)
				if xo.Size > 0 {
					_, ko := r.Header["Range"]
					switch ko {
					case true:
						ranges := strings.FieldsFunc(r.Header.Get("Range"), Split)

						if len(ranges) >= 2 {
							minrange, _ = strconv.Atoi(ranges[1])
							contentlenght = int(xo.Size) - minrange
						} else {
							contentlenght = int(xo.Size)
						}

						contentlenght = int(xo.Size) - minrange
						of = uint64(minrange)

						w.Header().Set("Content-Length", strconv.Itoa(contentlenght))
						w.Header().Set("Accept-Ranges", "bytes")
						w.Header().Set("Last-Modified", xo.ModTime.String())
						w.Header().Set("Content-Range", "bytes "+strconv.Itoa(minrange)+"-"+strconv.FormatUint(uint64(xo.Size-1), 10)+"/"+strconv.FormatUint(xo.Size, 10))
						w.Header().Set("Content-Type", mime)
						w.WriteHeader(http.StatusPartialContent)

						readFile(w, r, name, pool, xo, of)
						break
					case false:
						w.Header().Set("Content-Length", strconv.FormatUint(uint64(xo.Size), 10))
						readFile(w, r, name, pool, xo, of)
						break
					}
					break
				}

				fileparts = strings.Split(filename, ",")
				fileparts = fileparts[:len(fileparts)-1]
				for filepart := range fileparts {
					name = fileparts[filepart]
					xo, _ = ioctx.Stat(name)
					fsize = fsize + xo.Size
				}

				//fmt.Println("============== Req ==============")
				//for nnn, values := range r.Header {
				//	for _, value := range values {
				//		fmt.Println(nnn, value)
				//	}
				//}
				//fmt.Println("=================================")

				_, ko := r.Header["Range"]
				switch ko {
				case true:

					ranges := strings.FieldsFunc(r.Header.Get("Range"), Split)
					if len(ranges) >= 2 {
						minrange, _ = strconv.Atoi(ranges[1])
						contentlenght = int(fsize) - minrange
					} else {
						contentlenght = int(xo.Size)
					}

					sizes := []int{}
					actsz := []int{}
					before := 0

					for filepart := range fileparts {
						siz, _ := strconv.Atoi(strings.Split(fileparts[filepart], "-")[1])
						sizes = append(sizes, siz)
						x, ez := ioctx.Stat(fileparts[filepart])
						if ez != nil {
							wrados.Writelog("Can't get file info", fileparts[filepart])
						}
						actsz = append(actsz, int(x.Size))
					}

					for fp := range sizes {
						if minrange < sizes[fp] {
							for xd := range fileparts[:fp-1] { // Calculate prior file sizes
								before = actsz[xd] + before
								//fmt.Println(before)
							}
							fileparts = fileparts[fp-1:]
							sizes = sizes[fp-1:]
							break
						}
					}

					w.Header().Set("Content-Length", strconv.Itoa(contentlenght))
					w.Header().Set("Accept-Ranges", "bytes")
					w.Header().Set("Last-Modified", xo.ModTime.String())
					w.Header().Set("Content-Range", "bytes "+strconv.Itoa(minrange)+"-"+strconv.FormatUint(fsize-1, 10)+"/"+strconv.FormatUint(fsize, 10))
					w.Header().Set("Content-Type", mime)
					w.WriteHeader(http.StatusPartialContent)

					if minrange >= sizes[len(sizes)-1] {
						xo, _ = ioctx.Stat(name)
						of = xo.Size - uint64(contentlenght)
						_ = readFile(w, r, name, pool, xo, of)
					} else {
						for filepart := range fileparts {
							name = fileparts[filepart]
							xo, _ = ioctx.Stat(name)
							if filepart == 0 {
								of = uint64(minrange - before)

							} else {
								of = 0
							}
							x := readFile(w, r, name, pool, xo, of)
							if x == false {
								break
							}
						}

					}

				case false:
					w.Header().Set("Content-Length", strconv.FormatUint(fsize, 10))
					for filepart := range fileparts {
						name = fileparts[filepart]
						xo, _ = ioctx.Stat(name)
						x := readFile(w, r, name, pool, xo, 0)
						if x == false {
							break
						}
					}
				}
			}
		}
	}
}

func Geet(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]
	switch pool {
	case "favicon.ico":
		// DO nothing !
	default:
		name := strings.Join(s[2:], "/")
		if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
			randindex := rand.Intn(len(wrados.Rconnect.Connection))
			ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
			if e != nil {
				wrados.Writelog(e)
			}

			xo, lo := ioctx.Stat(name)
			//ss := "eeeeeeeeeeeeeee"
			ss, eror := metadata.DBClient(pool+"/"+name, "get", "")

			//filez := []string{}
			var filez []string

			//for n, values := range r.Header {
			//	if strings.HasPrefix(n, "Range") {
			//		for _, value := range values {
			//			fmt.Println(n)
			//			fmt.Println(value)
			//		}
			//	}
			//}

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
								//wrados.Writelog(er) // Broken Pipe, Connection Reset
								break
							}
							of = of + mx
							if of >= xo.Size {
								break
							}
						}
						wrados.Writelog(r.Method, xo.Size, "bytes", name, "from", pool)
					} else {
						errormsg := strings.Split(fmt.Sprint(lo), ",")
						wrados.Writelog(errormsg)
						msg := []byte(errormsg[len(errormsg)-1] + "\n")
						w.Header().Set("Content-Length", strconv.FormatUint(uint64(len(msg)), 10))
						switch errormsg[len(errormsg)-1] {
						case " No such file or directory":
							w.WriteHeader(http.StatusNotFound)
						default:
							w.WriteHeader(http.StatusInternalServerError)
						}
						_, _ = w.Write(msg)
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
							name := name + "-" + strconv.Itoa(size)
							xo, _ := ioct.Stat(name)
							_ = ioct.Write(name, mukuch, xo.Size)

							fileSegments = append(fileSegments, name)
							xxx = xxx + 1

							break
						}
						if len(mukuch) > configs.Conf.Uploadmaxpart {
							name := name + "-" + strconv.Itoa(size)
							xo, _ := ioct.Stat(name)
							_ = ioct.Write(name, mukuch, xo.Size)

							fileSegments = append(fileSegments, name)
							xxx = xxx + 1
							lenMukuch := len(mukuch)

							size = size + lenMukuch
							wrados.Writelog(r.Method, lenMukuch, "bytes, segment", name, "of", r.URL, "to", pool, size)
							mukuch = nil
						}
					}
					fileSegments = append(fileSegments, r.Header.Get("Content-Length"))
					wrados.Writelog("Created File", name, "In", pool)
					fmeta := strings.Join(fileSegments, ",")
					_, err := metadata.DBClient(pool+"/"+name, "set", fmeta)
					if err != nil {
						wrados.Writelog("error setting metadata:", err)
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

			ss, eror := metadata.DBClient(pool+"/"+name, "get", "")
			var filez []string

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
				_, _ = metadata.DBClient(pool+"/"+name, "del", "")

			}
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		msg := "Dangerous commands are disabled ! \n"
		wrados.Writelog(msg)
		_, _ = w.Write([]byte(msg))
	}
}

func Head(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]
	switch pool {
	case "favicon.ico":
		// DO nothing !
	default:
		name := strings.Join(s[2:], "/")
		if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
			randindex := rand.Intn(len(wrados.Rconnect.Connection))
			ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
			if e != nil {
				wrados.Writelog(e)
			}
			xo, lo := ioctx.Stat(name)
			ss, eror := metadata.DBClient(pool+"/"+name, "get", "")
			if eror != nil {
				wrados.Writelog(eror)
				break
			}

			if lo != nil {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(lo.Error()))
				_, _ = w.Write([]byte("\n"))
			} else {
				fmt.Println(" ")
				wrados.Writelog(xo, ss)
				//w.Header().Set("Content-Length", strconv.FormatUint(fsize, 10))
				for nnn, values := range r.Header {
					for _, value := range values {
						fmt.Println(nnn, value)
					}
				}
				wrados.Writelog(r.URL, r.Method, r.ContentLength, r.RequestURI)
				fmt.Println(" ")
			}
		} else {
			wrados.Writelog("Pool " + pool + " does not exists")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500 Internal Server Error \n"))
		}
	}
}
