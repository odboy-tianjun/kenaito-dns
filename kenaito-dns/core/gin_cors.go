package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               // 请求方法
		origin := c.Request.Header.Get("Origin") // 请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)  // 允许访问域
			c.Header("Access-Control-Allow-Methods", "POST") // 服务器支持的跨域请求的方法
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma, FooBar") // 允许跨域设置，可以返回其他子段
			c.Header("Access-Control-Max-Age", "86400")                                                                                                                                                                   // 24小时                                                                                                                                                                                                                                                                 // 缓存请求信息，单位为秒
			c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                          // 跨域请求是否需要带cookie信息，默认设置为true
			c.Set("content-type", "application/json;charset=utf8")                                                                                                                                                        // 设置返回格式是json
		}
		// 放行 OPTIONS 预检请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		c.Next()
	}
}
