package main

import (
	"configs"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	go web.PopulateUsers()
	go wrados.ListPools()
	web.RunServer()
}
