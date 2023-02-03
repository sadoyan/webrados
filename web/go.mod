module web

go 1.17

require (
	github.com/ceph/go-ceph v0.19.0
	configs v0.0.1
	metadata v0.0.1
	wrados v0.0.1
	auth v0.0.1
)

replace metadata => ../metadata
replace auth => ../auth
replace wrados => ../wrados
replace configs => ../configs
