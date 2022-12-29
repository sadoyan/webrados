module metadata

go 1.17

require (
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gomodule/redigo v1.8.9
)

require configs v0.0.1

require gopkg.in/ini.v1 v1.67.0 // indirect

replace configs => ../configs
