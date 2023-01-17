package main

import (
	"configs"
	"metadata"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	web.PopulatemMimes()
	go metadata.DBConnect()
	go web.PopulateUsers()
	go wrados.LsPools()
	go wrados.LogToFile()
	//go metadata.RedinitPool()
	//go metadata.MySQLInitPool()
	web.RunServer()
}
