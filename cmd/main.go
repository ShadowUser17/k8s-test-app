package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	var (
		router  = gin.New()
		address = flag.String("listen", ":8080", "Set address:port")
	)

	if config, err := rest.InClusterConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)

	} else {
		if clientSet, err := kubernetes.NewForConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)

		} else {
			router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
				return fmt.Sprintf(
					"[%s] %s -> %s %s %d %s\n",
					param.TimeStamp.Format(time.RFC1123),
					param.ClientIP,
					param.Method,
					param.Path,
					param.StatusCode,
					param.Latency,
				)
			}))

			router.GET("/", func(ctx *gin.Context) {
				if pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), meta.ListOptions{}); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})

				} else {
					var items = make([]string, 0)

					for index := range pods.Items {
						items = append(items, "pod/"+pods.Items[index].Name)
					}

					ctx.JSON(http.StatusOK, gin.H{"Pods": items})
				}
			})

			flag.Parse()
			if err := router.Run(*address); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
