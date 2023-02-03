module go-webrados

go 1.17

require (
	configs v0.0.1
	auth v0.0.1
	web v0.0.1
	wrados v0.0.1
)

require (
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/ceph/go-ceph v0.19.1-0.20230112054424-122159ed21a1 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace configs => ./configs
replace auth => ./auth
replace metadata => ./metadata
replace web => ./web
replace wrados => ./wrados
