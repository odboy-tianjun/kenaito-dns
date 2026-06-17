package dao

/*
 * @Description  解析配置定义与操作
 * @Author  https://www.odboy.cn
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
	var record ResolveVersion
	has, err := Engine.Table("resolve_version").
		Where("is_release = ?", 1).
		Desc("id").
		Limit(1).
		Get(&record)

	if err != nil || !has {
		return 0
	}
	return record.Version
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
	_, _ = Engine.Table("resolve_version").Update(ResolveVersion{IsRelease: 2}, ResolveVersion{})
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
	_, err := Engine.Table("resolve_version").Update(ResolveVersion{IsRelease: 2}, ResolveVersion{})
	if err != nil {
		return err
	}
	_, err = Engine.Table("resolve_version").Update(ResolveVersion{IsRelease: 1}, ResolveVersion{Version: version})
	if err != nil {
		return err
	}
	return nil
}

func FindResolveVersionPage(pageNo int, pageSize int) []*ResolveVersion {
	if pageSize <= 5 {
		pageSize = 5
	}
	if pageNo <= 0 {
		pageNo = 1
	}
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

// CountResolveVersionPage 返回版本总数（不分页）
func CountResolveVersionPage() int {
	session := Engine.Table("resolve_version")
	count, err := session.Count()
	if err != nil {
		fmt.Println(err)
	}
	return int(count)
}
