package main

/*
 * @Description  服务入口
 * @Author  www.odboy.cn
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
	"net/http"
	"time"
)

func main() {
	fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " kenaito-dns version = " + config.AppVersion)
	cache.ReloadCache()
	go initDNSServer()
	initRestfulServer()
}

func initDNSServer() {
	// 注册 DNS 请求处理函数
	dns.HandleFunc(".", core.HandleDNSRequest)
	// 设置服务器地址和协议
	server := &dns.Server{Addr: config.DNSServerPort, Net: "udp"}
	// 开始监听
	fmt.Printf("[dns]  [info]  "+time.Now().Format(config.AppTimeFormat)+" Starting DNS server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("[dns]  [error]  "+time.Now().Format(config.AppTimeFormat)+" Failed to start DNS server: %s\n", err.Error())
	}
}

func initRestfulServer() {
	gin.SetMode(config.WebMode)
	// 创建一个新的 Gin 引擎实例，使用 gin.New() 方法来创建一个不包含任何默认中间件的实例
	router := gin.New()
	// LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
	// 默认 gin.DefaultWriter = os.Stdout
	// 自定义的日志格式化函数
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return fmt.Sprintf("[gin]  [info]  %s [Request] [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			// 请求时间戳
			param.TimeStamp.Format(config.AppTimeFormat),
			// 客户端 IP 地址
			param.ClientIP,
			// 请求方法 (GET, POST 等)
			param.Method,
			// 请求路径
			param.Path,
			// 请求协议
			param.Request.Proto,
			// 响应状态码
			param.StatusCode,
			// 请求延迟时间
			param.Latency,
			// 用户代理
			param.Request.UserAgent(),
			// 错误信息（如果有的话）
			param.ErrorMessage,
		)
	}))
	// 允许使用跨域请求，全局中间件
	router.Use(core.Cors())
	// 使用 Recovery 中间件，处理任何出现的错误，并防止服务崩溃
	router.Use(gin.Recovery())
	server := &http.Server{
		Addr:         config.WebServerPort,
		Handler:      router,
		ReadTimeout:  config.WebReadTimeout * time.Second,
		WriteTimeout: config.WebWriteTimeout * time.Second,
	}
	controller.InitRestFunc(router)
	fmt.Printf("[gin]  [info]  "+time.Now().Format(config.AppTimeFormat)+" Start Gin server: %s\n", config.WebServerPort)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("[gin]  [error]  "+time.Now().Format(config.AppTimeFormat)+" Failed to start Gin server: %s\n", config.WebServerPort)
	}
}
