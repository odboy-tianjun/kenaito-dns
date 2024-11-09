package dao

/*
 * @Description  解析配置定义与操作
 * @Author  www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
)

type ResolveVersion struct {
	Id             int `xorm:"pk not null integer 'id' autoincr" json:"id"`
	CurrentVersion int `xorm:"not null integer 'curr_version'" json:"currentVersion"`
}

func (ResolveVersion) TableName() string {
	return "resolve_version"
}

func GetResolveVersion() int {
	var records []ResolveVersion
	err := Engine.Table("resolve_config").Where("`id` = ?", 1).Find(&records)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if len(records) == 0 {
		return 0
	}
	return records[0].CurrentVersion
}

func ModifyResolveVersion(currentVersion int) (bool, error) {
	wrapper := new(ResolveVersion)
	wrapper.Id = 1

	updateRecord := new(ResolveVersion)
	updateRecord.CurrentVersion = currentVersion
	_, err := Engine.Table("resolve_config").Update(updateRecord, wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
