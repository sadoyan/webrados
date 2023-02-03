package main

import (
	"auth"
	"configs"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	web.PopulatemMimes()
	go auth.PopulateBAusers()
	go wrados.LsPools()
	go wrados.LogToFile()
	web.RunServer()
}
