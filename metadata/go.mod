module metadata

go 1.17

require (
	configs v0.0.1
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/ceph/go-ceph v0.19.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	wrados v0.0.1
	tools v0.0.1
)

replace configs => ../configs

replace wrados => ../wrados
replace tools => ../tools

