module go-webrados

go 1.17

require (
	auth v0.0.1
	configs v0.0.1
	tools v0.0.1
	web v0.0.1
	wrados v0.0.1
)

require (
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/ceph/go-ceph v0.19.1-0.20230112054424-122159ed21a1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/itchyny/base58-go v0.2.0 // indirect
	github.com/leg100/surl v0.0.6-0.20221212214310-e628088f5204 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/sys v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	metadata v0.0.1 // indirect
)

replace configs => ./configs

replace auth => ./auth

replace metadata => ./metadata

replace web => ./web

replace wrados => ./wrados

replace tools => ./tools
