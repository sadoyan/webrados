module web

go 1.19

require github.com/ceph/go-ceph v0.19.0

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	configs v0.0.1
	metadata v0.0.1
	wrados v0.0.1

)

replace metadata => ../metadata

replace wrados => ../wrados

replace configs => ../configs
