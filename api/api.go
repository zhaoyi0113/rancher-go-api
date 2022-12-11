package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zhaoyi0113/rancher-go-api/internal"
)

func CreateRoute() *gin.Engine {
	var r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		fmt.Println("health check")
		fmt.Println(os.Environ())
		c.Writer.WriteHeader(http.StatusOK)
	})

	r.GET("/", func(c *gin.Context) {
		fmt.Println("health check from gcp")
		fmt.Println(os.Environ())
		c.Writer.WriteHeader(http.StatusOK)
	})

	r.POST("/transaction", func(c *gin.Context) {
		fmt.Println("get transation")
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Failed to read request", err)
			c.Writer.WriteHeader(http.StatusBadRequest)
		}
		fmt.Println("transaction", string(jsonData))
		internal.ProcessTransactionRequest(jsonData)
	})

	return r
}
