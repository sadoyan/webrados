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
	"github.com/leg100/surl"
	"io"
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
type adminapi struct {
	adminkey string
}
type signurl struct {
	signurl string
}

type credential struct {
	Keys    map[string]bool
	User    map[string]string
	Token   string
	SignURL string
	sync.RWMutex
}

var Credential = &credential{
	Keys:    map[string]bool{},
	User:    map[string]string{},
	Token:   "",
	SignURL: "",
	RWMutex: sync.RWMutex{},
}

func (ad *adminapi) auth() bool {
	if ad.adminkey == configs.Conf.AdminApiKey {
		return true
	} else {
		return false
	}
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
		return false

	} else {
		return true
	}

}
func (su *signurl) auth() bool {
	signer := surl.New([]byte(configs.Conf.JWTSecret))
	err := signer.Verify(su.signurl)
	if err != nil {
		//log.Println("verification failed: ", err.Error())
		return false
	}
	//fmt.Println("verification succeeded")
	return true
}

func DoAdminAuth(r *http.Request) bool {
	a := adminapi{adminkey: r.Header.Get("X-API-KEY")}
	return a.auth()
}

func DoAuth(r *http.Request) bool {
	switch {
	case configs.Conf.AuthApi:
		a := api{Key: r.Header.Get("X-API-KEY")}
		return a.auth()
	case configs.Conf.AuthBasic:
		username, password, ok := r.BasicAuth()
		if !ok {
			return false
		}
		b := basic{User: username, Pass: password, Auth: ok}
		return b.auth()
	case configs.Conf.AuthJWT:
		jwthdr, ok := r.URL.Query()["token"]
		if !ok {
			jwthdr = strings.Split(r.Header.Get("Authorization"), " ")
		}
		c := token{token: jwthdr[len(jwthdr)-1]}
		return c.auth()
	case configs.Conf.AuthSighn:
		vrfy := r.Host + r.RequestURI
		if r.TLS != nil {
			vrfy = "https://" + vrfy
		} else {
			vrfy = "http://" + vrfy
		}
		d := signurl{signurl: vrfy}
		return d.auth()

	}
	return false
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

type jwtinput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Expire   int    `json:"exp"`
}

func GenJWTtoken(in []byte) ([]byte, error) {
	var jwtin jwtinput
	err := json.Unmarshal(in, &jwtin)
	if err != nil {
		tools.WriteLogs(err)
		return nil, err
	}
	hmacSampleSecret := configs.Conf.JWTSecret
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(jwtin.Username+jwtin.Password)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
		"exp":  jwtin.Expire,
	})

	tokenString, err2 := token.SignedString(hmacSampleSecret)
	if err != nil {
		tools.WriteLogs("Error Getting JWT signed key:", err2)
		return nil, err2
	}
	return []byte(tokenString), nil
}

func SignUrl(body io.Reader) []byte {
	inputdata := map[string]int{}
	returnval := map[string]string{}

	k, _ := io.ReadAll(body)
	_ = json.Unmarshal(k, &inputdata)
	for z, x := range inputdata {
		signer := surl.New([]byte(configs.Conf.JWTSecret))
		dur := time.Duration(x) * time.Second
		signed, wer := signer.Sign(z, dur)
		if wer == nil {
			returnval[z] = signed
		}
	}
	j, _ := json.Marshal(returnval)
	return j
}
