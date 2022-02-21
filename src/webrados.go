package main

import (
	"configs"
	"metadata"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	go metadata.DBConnect()
	go web.PopulateUsers()
	go wrados.LsPools()
	go wrados.LogToFile()
	//go metadata.RedinitPool()
	//go metadata.MySQLInitPool()
	web.RunServer()
}
