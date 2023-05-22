module metadata

go 1.17

require (
	configs v0.0.1
	github.com/allegro/bigcache/v3 v3.1.0
	tools v0.0.1
	wrados v0.0.1
)

require (
	github.com/ceph/go-ceph v0.19.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace configs => ../configs

replace wrados => ../wrados

replace tools => ../tools
