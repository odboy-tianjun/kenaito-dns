package main

/*
 * @Description  服务入口
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"kenaito-dns/cache"
	"kenaito-dns/config"
	"kenaito-dns/controller"
	"kenaito-dns/core"
	"kenaito-dns/util"
	"net/http"
	"time"
)

func main() {
	fmt.Println("[app]  [info]  " + util.NowStr() + " kenaito-dns version = " + config.AppVersion)
	cache.ReloadCache()
	go initDNSServerUDP()
	go initDNSServerTCP()
	initRestfulServer()
}

func initDNSServerUDP() {
	dns.HandleFunc(".", core.HandleDNSRequest)
	server := &dns.Server{Addr: config.DNSServerPort, Net: "udp"}
	fmt.Printf("[dns]  [info]  "+util.NowStr()+" Starting DNS server on %s (UDP)\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("[dns]  [error]  "+util.NowStr()+" Failed to start DNS server (UDP): %s\n", err.Error())
	}
}

func initDNSServerTCP() {
	dns.HandleFunc(".", core.HandleDNSRequest)
	server := &dns.Server{Addr: config.DNSServerPort, Net: "tcp"}
	fmt.Printf("[dns]  [info]  "+util.NowStr()+" Starting DNS server on %s (TCP)\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("[dns]  [error]  "+util.NowStr()+" Failed to start DNS server (TCP): %s\n", err.Error())
	}
}

func initRestfulServer() {
	gin.SetMode(config.WebMode)
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[gin]  [info]  %s [Request] [%s] \\\"%s %s %s %d %s \\\"%s\\\" %s\\\"\\n",
			param.TimeStamp.Format(config.AppTimeFormat),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(core.Cors())
	router.Use(gin.Recovery())
	server := &http.Server{
		Addr:         config.WebServerPort,
		Handler:      router,
		ReadTimeout:  config.WebReadTimeout * time.Second,
		WriteTimeout: config.WebWriteTimeout * time.Second,
	}
	controller.InitRestFunc(router)
	fmt.Printf("[gin]  [info]  "+util.NowStr()+" Start Gin server: %s\n", config.WebServerPort)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("[gin]  [error]  "+util.NowStr()+" Failed to start Gin server: %s\n", err.Error())
	}
}
