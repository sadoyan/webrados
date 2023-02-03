module auth

go 1.17

require (
	configs v0.0.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	wrados v0.0.1
)

require (
	github.com/ceph/go-ceph v0.19.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace wrados => ../wrados

replace configs => ../configs
