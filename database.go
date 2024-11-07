package main

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Engine *xorm.Engine
)

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", "dns.sqlite3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("数据库引擎创建成功")

	err = Engine.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("数据库连接成功")
}
