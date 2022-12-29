module wrados

go 1.19

require configs v0.0.1

replace configs => ../configs

require github.com/ceph/go-ceph v0.19.0

require (
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
