module auth

go 1.17

require (
	configs v0.0.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/leg100/surl v0.0.6-0.20221212214310-e628088f5204
	tools v0.0.1
)

require (
	github.com/itchyny/base58-go v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace tools => ../tools

replace configs => ../configs
