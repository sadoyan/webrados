module web

go 1.17

require (
	auth v0.0.1
	configs v0.0.1
	github.com/ceph/go-ceph v0.19.0
	metadata v0.0.1
	tools v0.0.1
	wrados v0.0.1
)

require (
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace metadata => ../metadata

replace auth => ../auth

replace wrados => ../wrados

replace configs => ../configs

replace tools => ../tools
