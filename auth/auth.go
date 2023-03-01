package auth

import (
	"bufio"
	"configs"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"tools"
)

type auth interface {
	auth() bool
}

type basic struct {
	User string
	Pass string
	Auth bool
}
type token struct {
	token string
}
type api struct {
	Key string
}

type credential struct {
	Keys  map[string]bool
	User  map[string]string
	Token string
	sync.RWMutex
}

var Credential = &credential{
	Keys:    map[string]bool{},
	User:    map[string]string{},
	Token:   "",
	RWMutex: sync.RWMutex{},
}

func (ba *basic) auth() bool {

	md5HashInBytes := md5.Sum([]byte(ba.Pass))
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])
	pass, ok := Credential.User[ba.User]

	if !ok {
		return false
	}
	return md5HashInString == pass
}

func (ap *api) auth() bool {
	if _, ok := Credential.Keys[ap.Key]; ok {
		return true
	} else {
		return false
	}
}

func (tk *token) auth() bool {
	tok := tk.token
	hmacSampleSecret := configs.Conf.JWTSecret
	_, errr := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})
	if errr != nil {
		//tools.WriteLogs(tools.GetIP(r), r.Method, "JWT", errr, r.URL)
		return false

	} else {
		return true
	}

}

func authenticate(a auth) bool {
	return a.auth()
}

func AddUsers() {
	for {
		c, e := os.Open(configs.Conf.UsersFile)
		content := bufio.NewScanner(c)
		if e != nil {
			tools.WriteLogs(e)
		}

		for content.Scan() {
			if len(content.Text()) > 0 && !strings.HasPrefix(content.Text(), "#") {
				z := strings.Split(content.Text(), " ")
				if len(z) == 2 {
					if _, ok := Credential.User[z[0]]; !ok {
						// tools.WriteLogs("Found new user: " + z[0] + ", enabling!")
						Credential.Lock()
						Credential.User[z[0]] = z[1]
						Credential.Unlock()
					}
				} else if len(z) == 1 {
					if _, ok := Credential.Keys[z[0]]; !ok {
						// tools.WriteLogs("Found new apikey: " + z[0] + ", enabling!")
						Credential.Lock()
						Credential.Keys[z[0]] = true
						Credential.Unlock()
					}
				}
			}
		}
		_ = c.Close()
		//fmt.Println("-------------------")
		time.Sleep(10 * time.Second)
	}

}

func DoAuth(r *http.Request) bool {
	switch {
	case configs.Conf.AuthApi:
		a := api{Key: r.Header.Get("X-API-KEY")}
		if authenticate(&a) {
			return true
		} else {
			return false
		}
	case configs.Conf.AuthBasic:
		username, password, ok := r.BasicAuth()
		if !ok {
			return false
		}
		b := basic{User: username, Pass: password, Auth: ok}
		return authenticate(&b)
	case configs.Conf.AuthJWT:
		jwthdr := strings.Split(r.Header.Get("Authorization"), " ")
		c := token{token: jwthdr[len(jwthdr)-1]}
		return authenticate(&c)
	}
	return false
}
