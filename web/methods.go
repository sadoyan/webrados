package web

import (
	"bufio"
	"configs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"metadata"
	"net/http"
	"strconv"
	"strings"
	"wrados"

	"github.com/ceph/go-ceph/rados"
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

func Split(r rune) bool {
	return r == '=' || r == '-'
}

func readFile(w http.ResponseWriter, r *http.Request, name string, pool string, xo rados.ObjectStat, of uint64) bool {
	randindex := rand.Intn(len(wrados.Rconnect.Connection))
	ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)

	if e != nil {
		wrados.Writelog(e)
	}
	//mx := uint64(1024000)
	mx := uint64(256000)
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
	//start := time.Now()

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
				if len(filename) > 0 {
					var fileparts []string
					fileparts = strings.Split(filename, ",")
					lenq := fileparts[len(fileparts)-1]
					ff := fileparts[:len(fileparts)-1]
					w.Header().Set("Content-Length", lenq)
					for fp := range ff {
						name = fileparts[fp]
						xo, _ = ioctx.Stat(name)
						readFile(w, r, name, pool, xo, of)
					}

				} else {
					readFile(w, r, name, pool, xo, of)
				}

			case true:
				var fsize uint64
				var fileparts []string
				xo, _ = ioctx.Stat(name)
				if len(filename) == 0 {
					//if xo.Size > 0 {
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
	//elapsed := time.Since(start)
	//fmt.Println("Took :", elapsed)
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

				ssize := 2048000
				ssiz1 := 2048001

				if lenq < ssiz1 {
					reqBody, _ := ioutil.ReadAll(r.Body)
					_ = ioct.Write(name, reqBody, 0)
				} else {
					_ = ioct.Create(name, rados.CreateOption(lenq))
					fileSegments := make([]string, 0)
					reader := bufio.NewReader(r.Body)
					const BufferSize = 1
					buffer := make([]byte, BufferSize)
					writebuffer := make([]byte, 0)

					//start := time.Now()
					var bytecalc int
					var totalbytes int
					var segment string
					//PrintMemUsage()

					for {
						if bytecalc == 0 {
							if lenq < configs.Conf.Uploadmaxpart-ssiz1 {
								segment = name
							} else {
								segment = name + "-0"
								fileSegments = append(fileSegments, segment)
								wrados.Writelog(r.Method, bytecalc, "bytes, segment", segment, "of", r.URL, "to", pool, totalbytes)
							}
						}
						_, eerr := reader.Read(buffer)
						if eerr != nil {
							_ = ioct.Append(segment, writebuffer)
							wrados.Writelog(r.Method, bytecalc, "bytes, segment", segment, "of", r.URL, "to", pool, totalbytes)
							writebuffer = nil
							break
						}

						writebuffer = append(writebuffer, buffer...)
						if bytecalc >= configs.Conf.Uploadmaxpart-ssiz1 {
							segment = name + "-" + strconv.Itoa(totalbytes)
							_ = ioct.Create(segment, rados.CreateOption(configs.Conf.Uploadmaxpart-ssiz1))
							fileSegments = append(fileSegments, segment)
							wrados.Writelog(r.Method, bytecalc, "bytes, segment", segment, "of", r.URL, "to", pool, totalbytes)
							bytecalc = 0
						}
						if len(writebuffer) == ssize {
							_ = ioct.Append(segment, writebuffer)
							writebuffer = nil
						}
						bytecalc = bytecalc + 1
						totalbytes = totalbytes + 1
					}
					//PrintMemUsage()
					if lenq >= configs.Conf.Uploadmaxpart-ssiz1 {
						fileSegments = append(fileSegments, r.Header.Get("Content-Length"))
						fmeta := strings.Join(fileSegments, ",")
						_, metaerr := metadata.DBClient(pool+"/"+name, "set", fmeta)
						if metaerr != nil {
							wrados.Writelog("error setting metadata:", metaerr)
						}
					}

					wrados.Writelog("Created File", name, "In", pool)
					//log.Printf("Execution time %s\n", time.Since(start))
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
