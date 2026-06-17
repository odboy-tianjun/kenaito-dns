package dao

/*
 * @Description  连接数据库
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"kenaito-dns/config"
	"kenaito-dns/util"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	Engine *xorm.Engine
)

func init() {
	var err error
	Engine, err = xorm.NewEngine(config.DataSourceDriverName, config.DataSourceName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[xorm] [info]  " + util.NowStr() + " 数据库引擎创建成功(database engine create success)")
	// 连接池配置
	Engine.SetMaxOpenConns(config.DataSourceMaxOpenConnectionSize)
	Engine.SetMaxIdleConns(config.DataSourceMaxIdleConnectionSize)
	Engine.SetConnMaxLifetime(config.DataSourceConnMaxLifetime * time.Minute)
	// 日志相关
	Engine.ShowSQL(config.DataSourceShowLog)
	Engine.Logger().SetLevel(config.DataSourceLogLevel)
	// 自动建表
	err = Engine.Sync(new(ResolveRecord), new(ResolveVersion))
	if err != nil {
		fmt.Printf("[xorm] [error]  %s 数据库同步表结构失败: %v\n", util.NowStr(), err)
	} else {
		fmt.Println("[xorm] [info]  " + util.NowStr() + " 数据库表结构同步成功(database schema synced)")
	}
	err = Engine.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[xorm] [info]  " + util.NowStr() + " 数据库连接成功(database engine connect success)")
}
