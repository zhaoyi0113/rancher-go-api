package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoute() *gin.Engine {
	var r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		fmt.Println("health check")
		c.Writer.WriteHeader(http.StatusOK)
	})

	return r
}
