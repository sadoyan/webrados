package main

import (
	"configs"
	"web"
	"wrados"
)

func main() {
	configs.SetVarsik()
	wrados.ListPools()
	web.RunServer()
}
