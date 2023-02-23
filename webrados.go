package main

import (
	"auth"
	"configs"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	web.HttpMimes.Populate()
	go auth.PopulateUsers()
	go wrados.LsPools()
	go wrados.LogToFile()
	web.RunServer()
}
