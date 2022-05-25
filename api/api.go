package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoute() *gin.Engine {
	var r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		fmt.Println("health check")
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
	})

	return r
}
