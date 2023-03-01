module auth

go 1.17

require (
	configs v0.0.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	tools v0.0.1
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

replace tools => ../tools

replace configs => ../configs
