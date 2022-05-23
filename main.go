package main

import "github.com/zhaoyi0113/rancher-go-api/api"

func main() {
	route := api.CreateRoute()
	route.Run()
}
