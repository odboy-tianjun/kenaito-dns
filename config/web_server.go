package config

/*
 * @Description  Web服务配置
 * @Author  www.odboy.cn
 * @Date  20241108
 */
import "github.com/gin-gonic/gin"

const (
	WebServerPort   = ":18001"
	WebMode         = gin.ReleaseMode
	WebReadTimeout  = 10
	WebWriteTimeout = 10
)
