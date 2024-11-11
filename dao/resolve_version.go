package dao

/*
 * @Description  解析配置定义与操作
 * @Author  www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"kenaito-dns/config"
	"time"
)

type ResolveVersion struct {
	Id         int    `xorm:"pk not null integer 'id' autoincr" json:"id"`
	Version    int    `xorm:"not null integer 'version'" json:"version"`
	CreateTime string `xorm:"not null text 'create_time'" json:"createTime"`
	IsRelease  int    `xorm:"not null integer 'is_release'" json:"isRelease"`
}

func (ResolveVersion) TableName() string {
	return "resolve_version"
}

func GetResolveVersion() int {
	var records []ResolveVersion
	session := Engine.Table("resolve_version")
	session.Desc("id")
	session.And("is_release = ?", 1)
	err := session.Find(&records)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if len(records) == 0 {
		return 0
	}
	return records[0].Version
}

func FindResolveVersionByVersion(version int) *ResolveVersion {
	var record ResolveVersion
	session := Engine.Table("resolve_version")
	session.And("version = ?", version)
	result, err := session.Get(&record)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if result {
		return &record
	}
	return nil
}

func SaveResolveVersion(wrapper *ResolveVersion) (bool, error) {
	// 全表更新为未发布
	_, _ = Engine.Table("resolve_version").Update(ResolveVersion{IsRelease: 2}, ResolveVersion{})
	// 新增发布记录
	wrapper.CreateTime = time.Now().Format(config.DataTimeFormat)
	wrapper.IsRelease = 1
	_, err := Engine.Table("resolve_version").Insert(wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func ModifyResolveVersion(version int) error {
	// 全表更新为未发布
	_, err := Engine.Table("resolve_version").Update(ResolveVersion{IsRelease: 2}, ResolveVersion{})
	if err != nil {
		return err
	}
	// 更新具体的版本为已发布
	_, err = Engine.Table("resolve_version").Update(ResolveVersion{IsRelease: 1}, ResolveVersion{Version: version})
	if err != nil {
		return err
	}
	return nil
}

func FindResolveVersionPage(pageNo int, pageSize int) []*ResolveVersion {
	// 每页显示5条记录
	if pageSize <= 5 {
		pageSize = 5
	}
	// 要查询的页码
	if pageNo <= 0 {
		pageNo = 1
	}
	// 计算跳过的记录数
	offset := (pageNo - 1) * pageSize
	records := make([]*ResolveVersion, 0)
	session := Engine.Table("resolve_version")
	session.Desc("id")
	err := session.Limit(pageSize, offset).Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}
func CountResolveVersionPage(pageNo int, pageSize int) int {
	// 每页显示5条记录
	if pageSize <= 5 {
		pageSize = 5
	}
	// 要查询的页码
	if pageNo <= 0 {
		pageNo = 1
	}
	// 计算跳过的记录数
	offset := (pageNo - 1) * pageSize
	session := Engine.Table("resolve_version")
	count, err := session.Limit(pageSize, offset).Count()
	if err != nil {
		fmt.Println(err)
	}
	return int(count)
}
