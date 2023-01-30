module metadata

go 1.17

require (
	configs v0.0.1
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/ceph/go-ceph v0.19.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gomodule/redigo v1.8.9
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	wrados v0.0.1
)

replace configs => ../configs

replace wrados => ../wrados
