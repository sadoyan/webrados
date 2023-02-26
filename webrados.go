package main

import (
	"auth"
	"configs"
	"tools"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	web.HttpMimes.Populate()
	go auth.PopulateUsers()
	go wrados.LsPools()
	go tools.LogToFile()
	web.RunServer()
}
