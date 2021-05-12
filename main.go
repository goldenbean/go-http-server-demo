package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	httpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "http request count",
		},
		[]string{"endpoint"})

	httpRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration",
			Help: "http request duration",
		},
		[]string{"endpoint"})

	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "yarn_metrics_http_ops_total",
		Help: "The total number of processed events",
	})

	homePath string
)

func init() {
	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestDuration)

	loadConfig()
}

func loadConfig() {
	config := viper.New()

	config.SetConfigName("config") // name of config file (without extension)
	config.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath(".")      // optionally look for config in the working directory
	err := config.ReadInConfig()   // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	homePath = config.GetString("home")

}

func server() {

	router := gin.Default()

	router.GET("/ok.htm", func(c *gin.Context) {
		// opsProcessed.Inc()

		ret := "gin, hello world, from "

		addrs, err := net.InterfaceAddrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			fmt.Println(addr.String())
			ret += addr.String() + "; "
		}

		ret += "\n"
		sec := time.Now().Unix()
		ret += "timestamp: "
		ret += strconv.FormatInt(int64(sec), 10)

		c.String(http.StatusOK, ret)
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Static("/static", "./static")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "static/gitbook/gitbook/index.html", nil)
	// })
	router.Use(static.Serve("/", static.LocalFile(homePath, true)))
	//router.Static("/gitbook", "./static/gitbook")

	router.Run(":8080")
}

func main() {

	server()
}
