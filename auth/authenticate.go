package auth

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"
	"wrados"
)

var users = map[string]string{
	//"test":  "secret",
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
