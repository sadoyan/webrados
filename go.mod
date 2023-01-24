module go-webrados

go 1.17

require (
	configs v0.0.1
	metadata v0.0.1
	web v0.0.1
	wrados v0.0.1
)

require (
	github.com/ceph/go-ceph v0.19.1-0.20230112054424-122159ed21a1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	golang.org/x/sys v0.4.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace configs => ./configs

replace metadata => ./metadata

replace web => ./web

replace wrados => ./wrados
