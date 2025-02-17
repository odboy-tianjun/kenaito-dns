package config

/*
 * @Description  数据库配置
 * @Author  https://www.odboy.cn
 * @Date  20241108
 */
import "xorm.io/core"

const (
	DataSourceDriverName            = "sqlite3"      // 驱动名称
	DataSourceName                  = "dns.sqlite3"  // 数据源名称
	DataSourceMaxOpenConnectionSize = 30             // 最大db连接数
	DataSourceMaxIdleConnectionSize = 10             // 最大db连接空闲数
	DataSourceConnMaxLifetime       = 30             // 超过空闲数连接存活时间
	DataSourceShowLog               = true           // 是否显示SQL语句
	DataSourceLogLevel              = core.LOG_DEBUG // 日志级别
)
