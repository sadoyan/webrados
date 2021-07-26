package main

import (
	"configs"
	"metadata"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	go web.PopulateUsers()
	go wrados.LsPools()
	go wrados.LogToFile()
	go metadata.RedinitPool()
	web.RunServer()
}
