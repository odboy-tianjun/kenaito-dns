package dao

/*
 * @Description  连接数据库
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"kenaito-dns/config"
	"kenaito-dns/util"
	"time"
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
	err = Engine.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[xorm] [info]  " + util.NowStr() + " 数据库连接成功(database engine connect success)")
}
