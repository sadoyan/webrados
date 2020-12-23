package main

import (
	"configs"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	go wrados.ListPools()
	web.RunServer()
}
