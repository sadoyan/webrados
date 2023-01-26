module metadata

go 1.17

require (
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gomodule/redigo v1.8.9
	configs v0.0.1
    wrados v0.0.1
	github.com/ceph/go-ceph v0.19.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace configs => ../configs
replace wrados => ../wrados
