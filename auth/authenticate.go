package auth

import (
	"bufio"
	"configs"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
	"wrados"
)

var BAusers = map[string]string{
	//"test":  "secret",
}

func PopulateBAusers() {
	for {
		c, e := os.Open(configs.Conf.UsersFile)
		content := bufio.NewScanner(c)
		if e != nil {
			wrados.Writelog(e)
		}

		for content.Scan() {
			if len(content.Text()) > 0 {
				z := strings.Split(content.Text(), " ")
				if _, ok := BAusers[z[0]]; !ok {
					wrados.Writelog("Found new user: " + z[0] + ", enabling!")
					BAusers[z[0]] = z[1]
				}
			}
		}
		_ = c.Close()

		time.Sleep(10 * time.Second)
	}

}

func isBAauthorised(username, password string) bool {

	md5HashInBytes := md5.Sum([]byte(password))
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])

	pass, ok := BAusers[username]

	if !ok {
		return false
	}
	return md5HashInString == pass
}

type jwtinput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenJWTtoken(in []byte) ([]byte, error) {
	var jwtin jwtinput
	err := json.Unmarshal(in, &jwtin)
	if err != nil {
		wrados.Writelog(err)
		return nil, err
	}

	hmacSampleSecret := configs.Conf.JWTSecret
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"username": jwtin.Username,
	//	"password": jwtin.Password,
	//})
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(jwtin.Username+jwtin.Password)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
	})

	tokenString, err2 := token.SignedString(hmacSampleSecret)
	if err != nil {
		wrados.Writelog("Error Getting JWT signed key:", err2)
		return nil, err2
	}
	return []byte(tokenString), nil
}

func CheckJWTtoken(tok string, r *http.Request) bool {
	hmacSampleSecret := configs.Conf.JWTSecret
	_, errr := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})
	if errr != nil {
		wrados.Writelog(configs.GetIP(r), r.Method, "JWT", errr, r.URL)
		return false

	} else {
		return true
	}

}

func CheckAuth(w http.ResponseWriter, r *http.Request) bool {
	const unauth = http.StatusUnauthorized

	switch {
	case configs.Conf.AuthApi:
		if r.Header.Get("X-API-KEY") == configs.Conf.Apikey {
			return true
		} else {
			wrados.Writelog(configs.GetIP(r), r.Method, "Invalid APIKEY", r.URL)
			http.Error(w, http.StatusText(unauth), unauth)
			return false
		}
	case configs.Conf.AuthJWT:
		jwthdr := strings.Split(r.Header.Get("Authorization"), " ")
		if CheckJWTtoken(jwthdr[len(jwthdr)-1], r) {
			return true
		} else {
			http.Error(w, http.StatusText(unauth), unauth)
			return false
		}
	case configs.Conf.AuthBasic:
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(unauth), unauth)
			wrados.Writelog(configs.GetIP(r), r.Method, "401 Unauthorized: No basic auth present", r.URL)
			return false
		}

		if !isBAauthorised(username, password) {
			http.Error(w, http.StatusText(unauth), unauth)
			wrados.Writelog(configs.GetIP(r), r.Method, "401 Unauthorized: Invalid Credentials", r.URL)
			return false
		}
		w.WriteHeader(http.StatusOK)
		return true
	default:
		return false
	}
}
