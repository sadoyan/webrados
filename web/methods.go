package web

import (
	"bufio"
	"configs"
	"encoding/json"
	"fmt"
	"github.com/ceph/go-ceph/rados"
	"io/ioutil"
	"math/rand"
	"metadata"
	"net/http"
	"strconv"
	"strings"
	"tools"
	"wrados"
)

var minrange int
var contentlenght int
var of uint64

func respCodewriter(f error, w http.ResponseWriter, r *http.Request) string {
	if strings.Split(f.Error(), ",")[1] == " No such file or directory" {
		w.WriteHeader(http.StatusNotFound)
		tools.WriteLogs(tools.GetIP(r), r.Method, f.Error(), r.URL.String())
		return http.StatusText(404) + ": " + r.URL.String() + "\n"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		tools.WriteLogs(tools.GetIP(r), r.Method, f.Error(), r.URL.String())
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
		tools.WriteLogs("Error opening", e)
		return false
	}

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
			tools.WriteLogs("Error reading segment", name, err)
			break
		}
		_, er := w.Write(bytesOut)

		if er != nil {
			if !strings.HasPrefix(er.Error(), "write tcp") {
				tools.WriteLogs("Broken pipe", er)
			}
			return false
		}
		of = of + mx
		if of >= xo.Size {
			break
		}
	}
	tools.WriteLogs(tools.GetIP(r), r.Method, xo.Size, "bytes", name, "from", pool)
	return true
}

func Get(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	pool := s[1]
	if pool == "favicon.ico" {
		return
	}
	name := strings.Join(s[len(s)-1:], "/")
	extension := strings.Split(name, ".")[1]
	_, ok := wrados.Rconnect.Poolnames[pool]
	if !ok {
		tools.WriteLogs("Error connecting to pool", pool)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 Internal Server Error \n"))
		return
	}
	randindex := rand.Intn(len(wrados.Rconnect.Connection))
	ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
	defer ioctx.Destroy()
	if e != nil {
		tools.WriteLogs(e)
	}
	filename, eror := metadata.DBClient(pool+"/"+name, "get", "")
	xo, lo := ioctx.Stat(name)

	if lo != nil {
		errormsg := strings.Split(lo.Error(), ",")

		tools.WriteLogs(tools.GetIP(r), r.Method, r.URL, errormsg[len(errormsg)-1])
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
				var size int
				if len(ns) > 1 {
					size, _ = strconv.Atoi(ns[len(ns)-1])
				} else {
					size = int(xo.Size)
				}
				fileInfo := &FileInfo{
					Size:  size,
					Pool:  pool,
					Parts: len(ns),
					Name:  name,
				}
				//b, _ := json.Marshal(fileInfo)
				b, _ := json.MarshalIndent(fileInfo, "", "    ")
				_, _ = w.Write(b)
				_, _ = w.Write([]byte("\n"))
			}
		}
		return
	}

	mime, mok := HttpMimes.Lookup(extension)
	w.Header().Set("Last-Modified", xo.ModTime.String())
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
			xo, _ = ioctx.Stat(name)
			w.Header().Set("Content-Length", strconv.FormatUint(xo.Size, 10))
			readFile(w, r, name, pool, xo, of)
		}
	case true:
		var fsize uint64
		var fileparts []string

		xoSizeInt := int(xo.Size)
		xoSizeStr := strconv.Itoa(xoSizeInt)

		w.Header().Set("Content-Length", strconv.Itoa(contentlenght))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Last-Modified", xo.ModTime.String())
		w.Header().Set("Content-Type", mime)

		if len(filename) == 0 {
			_, ko := r.Header["Range"]
			switch ko {
			case true:
				ranges := strings.FieldsFunc(r.Header.Get("Range"), Split)

				if len(ranges) >= 2 {
					minrange, _ = strconv.Atoi(ranges[1])
					contentlenght = xoSizeInt - minrange
				} else {
					contentlenght = xoSizeInt
				}

				contentlenght = xoSizeInt - minrange
				of = uint64(minrange)

				w.Header().Set("Content-Range", "bytes "+strconv.Itoa(minrange)+"-"+strconv.Itoa(xoSizeInt-1)+"/"+xoSizeStr)
				w.WriteHeader(http.StatusPartialContent)

				readFile(w, r, name, pool, xo, of)
				break
			case false:
				w.Header().Set("Content-Length", xoSizeStr)
				readFile(w, r, name, pool, xo, of)
				break
			}
			break
		}

		fileparts = strings.Split(filename, ",")
		fsize, _ = strconv.ParseUint(fileparts[len(fileparts)-1], 10, 64)
		fileparts = fileparts[:len(fileparts)-1]

		//for _, filepart := range fileparts {
		//	xo, _ = ioctx.Stat(filepart)
		//	fsize = fsize + xo.Size
		//}

		_, ko := r.Header["Range"]
		switch ko {
		case true:
			ranges := strings.FieldsFunc(r.Header.Get("Range"), Split)

			if len(ranges) >= 2 {
				minrange, _ = strconv.Atoi(ranges[1])
				contentlenght = int(fsize) - minrange
			} else {
				contentlenght = xoSizeInt
			}

			sizes := []int{}
			actsz := []int{}
			before := 0

			for _, filepart := range fileparts {
				siz, _ := strconv.Atoi(strings.Split(filepart, "-")[1])
				sizes = append(sizes, siz)
				x, ez := ioctx.Stat(filepart)
				if ez != nil {
					tools.WriteLogs("Can't get file info", filepart)
				}
				actsz = append(actsz, int(x.Size))
			}

			for fp := range sizes {
				if minrange < sizes[fp] {
					for xd := range fileparts[:fp-1] { // Calculate prior file sizes
						before = actsz[xd] + before
					}
					fileparts = fileparts[fp-1:]
					sizes = sizes[fp-1:]
					break
				}
			}

			w.Header().Set("Content-Range", "bytes "+strconv.Itoa(minrange)+"-"+strconv.FormatUint(fsize-1, 10)+"/"+strconv.FormatUint(fsize, 10))
			w.WriteHeader(http.StatusPartialContent)

			if minrange >= sizes[len(sizes)-1] {
				filepart := fileparts[len(fileparts)-1]
				xo, _ = ioctx.Stat(filepart)
				of = xo.Size - uint64(contentlenght)
				_ = readFile(w, r, filepart, pool, xo, of)
			} else {
				for f, filepart := range fileparts {
					xo, _ = ioctx.Stat(filepart)
					if f == 0 {
						of = uint64(minrange - before)
					} else {
						of = 0
					}
					x := readFile(w, r, filepart, pool, xo, of)
					if x == false {
						break
					}
				}

			}
		case false:
			w.Header().Set("Content-Length", strconv.FormatUint(fsize, 10))
			for _, filepart := range fileparts {
				xo, _ = ioctx.Stat(filepart)
				x := readFile(w, r, filepart, pool, xo, 0)
				if x == false {
					break
				}
			}
		}
	}
}

func Put(w http.ResponseWriter, r *http.Request) {

	if configs.Conf.Readonly {
		msg := "Server is running in read only mode !"
		tools.WriteLogs(tools.GetIP(r), msg)
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(msg + "\n"))
		return
	}
	s := strings.Split(r.URL.Path, "/")
	if len(s) < 3 {
		tools.WriteLogs(tools.GetIP(r), "Invalid pool name")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500: Invalid pool name \n"))
		return
	}
	pool := s[1]
	name := strings.Join(s[2:], "/")
	if _, ok := wrados.Rconnect.Poolnames[pool]; !ok {
		tools.WriteLogs(tools.GetIP(r), "Pool not found")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("404: Not Found \n"))
		return
	}
	randindex := rand.Intn(len(wrados.Rconnect.Connection))
	ioct, _ := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
	defer ioct.Destroy()
	lenq, lqe := strconv.Atoi(r.Header.Get("Content-Length"))
	if lqe != nil {
		tools.WriteLogs(tools.GetIP(r), "Invalid pool name")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("403: Content-Length header is mandatory\n"))
		return
	}
	ssize := 2048000
	ssiz1 := 2048001
	switch lenq < ssiz1 {
	case true:
		reqBody, _ := ioutil.ReadAll(r.Body)
		_ = ioct.Write(name, reqBody, 0)
	case false:
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
					tools.WriteLogs(tools.GetIP(r), r.Method, bytecalc, "bytes, segment", segment, "of", r.URL, "to", pool, totalbytes)
				}
			}
			_, eerr := reader.Read(buffer)
			if eerr != nil {
				_ = ioct.Append(segment, writebuffer)
				tools.WriteLogs(tools.GetIP(r), r.Method, bytecalc, "bytes, segment", segment, "of", r.URL, "to", pool, totalbytes)
				writebuffer = nil
				break
			}

			writebuffer = append(writebuffer, buffer...)
			if bytecalc >= configs.Conf.Uploadmaxpart-ssiz1 {
				segment = name + "-" + strconv.Itoa(totalbytes)
				_ = ioct.Create(segment, rados.CreateOption(configs.Conf.Uploadmaxpart-ssiz1))
				fileSegments = append(fileSegments, segment)
				tools.WriteLogs(tools.GetIP(r), r.Method, bytecalc, "bytes, segment", segment, "of", r.URL, "to", pool, totalbytes)
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
				tools.WriteLogs("error setting metadata:", metaerr)
			}
			_ = ioct.Append(name, []byte(fmeta))
		}

		tools.WriteLogs("Created File", name, "In", pool)
		//log.Printf("Execution time %s\n", time.Since(start))
	}
	tools.WriteLogs(tools.GetIP(r), r.Method, r.Header.Get("Content-Length"), "bytes", r.URL, "to", pool)
}

func Del(w http.ResponseWriter, r *http.Request) {
	switch configs.Conf.DangeZone {
	case true:
		s := strings.Split(r.URL.Path, "/")

		if len(s) >= 3 {
			pool := s[1]
			name := strings.Join(s[2:], "/")

			_, delcache := r.URL.Query()["cache"]
			if delcache {
				_, _ = metadata.DBClient(pool+"/"+name, "del", "")
				tools.WriteLogs(tools.GetIP(r), r.Method, "removing", pool+"/"+name, "from cache")
				return
			}
			ss, eror := metadata.DBClient(pool+"/"+name, "get", "")
			var filez []string
			if _, ok := wrados.Rconnect.Poolnames[pool]; ok {
				randindex := rand.Intn(len(wrados.Rconnect.Connection))
				ioct, _ := wrados.Rconnect.Connection[randindex].OpenIOContext(pool)
				defer ioct.Destroy()

				filez = append(filez, name)
				if eror == nil {
					//filez = append(filez, name)

					fileparts := strings.Split(ss, ",")
					fileparts = fileparts[:len(fileparts)-1]
					for _, filepart := range fileparts {
						filez = append(filez, filepart)
					}

				}

				for _, filename := range filez {
					f := ioct.Delete(filename)
					if f != nil {
						_, _ = fmt.Fprintf(w, respCodewriter(f, w, r))
					} else {
						tools.WriteLogs(tools.GetIP(r), r.Method, filename, "from", pool)
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
		tools.WriteLogs(tools.GetIP(r), msg)
		_, _ = w.Write([]byte(msg))
	}
}
