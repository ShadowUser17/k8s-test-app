package main

import (
	_ "context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/rest"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	var (
		router  = gin.New()
		address = flag.String("listen", ":8080", "Set address:port")
	)

	router.Use(gin.LoggerWithFormatter(LoggerFormatter))
	router.GET("/", HandleSecrets)

	flag.Parse()
	if err := router.Run(*address); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func HandleSecrets(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Status": "OK"})
}

func LoggerFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf(
		"[%s] %s -> %s %s %d %s\n",
		param.TimeStamp.Format(time.RFC1123),
		param.ClientIP,
		param.Method,
		param.Path,
		param.StatusCode,
		param.Latency,
	)
}
