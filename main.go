package main

import (
	"fmt"
	"os"

	"github.com/zhaoyi0113/rancher-go-api/api"
)

func main() {
	fmt.Println(os.Environ())
	route := api.CreateRoute()
	route.Run()
}
