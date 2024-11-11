package dao

/*
 * @Description  解析记录定义与操作
 * @Author  www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"kenaito-dns/config"
	"kenaito-dns/domain"
	"kenaito-dns/util"
	"strings"
	"time"
)

type ResolveRecord struct {
	Id         int    `xorm:"pk not null integer 'id' autoincr" json:"id"`
	Name       string `xorm:"not null text 'name'" json:"name"`
	RecordType string `xorm:"not null text 'record_type'" json:"recordType"`
	Ttl        int    `xorm:"not null integer 'ttl'" json:"ttl"`
	Value      string `xorm:"not null text 'value'" json:"value"`
	Version    int    `xorm:"not null integer 'version'" json:"version"`
	CreateTime string `xorm:"not null text 'create_time'" json:"createTime"`
	UpdateTime string `xorm:"not null text 'update_time'" json:"updateTime"`
	Enabled    int    `xorm:"not null integer 'enabled'" json:"enabled"`
}

func (ResolveRecord) TableName() string {
	return "resolve_record"
}

func FindResolveRecordById(id int) *ResolveRecord {
	var record ResolveRecord
	_, err := Engine.Table("resolve_record").Where("`id` = ?", id).Get(&record)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &record
}

func FindOneResolveRecord(wrapper *ResolveRecord, version int) *ResolveRecord {
	var record ResolveRecord
	_, err := Engine.Table("resolve_record").Where("`name` = ? and `record_type` = ? and `value` = ? and `version` = ?",
		wrapper.Name,
		wrapper.RecordType,
		wrapper.Value,
		version,
	).Get(&record)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &record
}

func FindResolveRecordByVersion(version int, isAll bool) []ResolveRecord {
	var records []ResolveRecord
	session := Engine.Table("resolve_record")
	session.Where("`version` = ?", version)
	if !isAll {
		session.Where("`enabled` = ?", 1)
	}
	err := session.Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}

func FindResolveRecordByNameType(name string, recordType string) []ResolveRecord {
	var records []ResolveRecord
	err := Engine.Table("resolve_record").Where("`name` = ? and `record_type` = ? and `version` = ?", name, recordType, GetResolveVersion()).Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}

func FindResolveRecordPage(pageNo int, pageSize int, args *domain.QueryPageArgs) []*ResolveRecord {
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
	records := make([]*ResolveRecord, 0)
	session := Engine.Table("resolve_record").Where("")
	if args != nil {
		if !util.IsBlank(args.Name) {
			qs := "%" + strings.TrimSpace(args.Name) + "%"
			session.And("`name` LIKE ?", qs)
		}
		if !util.IsBlank(args.Type) {
			qs := strings.TrimSpace(args.Type)
			session.And("`record_type` = ?", qs)
		}
		if !util.IsBlank(args.Value) {
			qs := strings.TrimSpace(args.Value)
			session.And("`value` = ?", qs)
		}
	}
	session.And("`version` = ?", GetResolveVersion())
	err := session.Limit(pageSize, offset).Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}
func CountResolveRecordPage(pageNo int, pageSize int, args *domain.QueryPageArgs) int {
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
	session := Engine.Table("resolve_record").Where("")
	if args != nil {
		if !util.IsBlank(args.Name) {
			qs := "%" + strings.TrimSpace(args.Name) + "%"
			session.And("`name` LIKE ?", qs)
		}
		if !util.IsBlank(args.Type) {
			qs := strings.TrimSpace(args.Type)
			session.And("`record_type` = ?", qs)
		}
		if !util.IsBlank(args.Value) {
			qs := strings.TrimSpace(args.Value)
			session.And("`value` = ?", qs)
		}
	}
	session.And("`version` = ?", GetResolveVersion())
	count, err := session.Limit(pageSize, offset).Count()
	if err != nil {
		fmt.Println(err)
	}
	return int(count)
}

func SaveResolveRecord(wrapper *ResolveRecord) (bool, error) {
	wrapper.CreateTime = time.Now().Format(config.DataTimeFormat)
	wrapper.UpdateTime = time.Now().Format(config.DataTimeFormat)
	wrapper.Enabled = 1
	_, err := Engine.Table("resolve_record").Insert(wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func BackupResolveRecord(record *ResolveRecord) (bool, error, int, int) {
	var backupRecords []*ResolveRecord
	oldVersion := GetResolveVersion()
	newVersion := GetResolveVersion() + 1
	oldRecords := FindResolveRecordByVersion(oldVersion, true)
	for _, oldRecord := range oldRecords {
		newRecord := new(ResolveRecord)
		newRecord.Name = oldRecord.Name
		newRecord.RecordType = oldRecord.RecordType
		newRecord.Ttl = oldRecord.Ttl
		newRecord.Value = oldRecord.Value
		newRecord.Version = newVersion
		newRecord.CreateTime = oldRecord.CreateTime
		newRecord.UpdateTime = oldRecord.UpdateTime
		newRecord.Enabled = oldRecord.Enabled
		backupRecords = append(backupRecords, newRecord)
	}
	record.Version = newVersion
	if len(backupRecords) > 0 {
		_, err := Engine.Table("resolve_record").Insert(backupRecords)
		if err != nil {
			return false, err, 0, 0
		}
	}
	updRecord := new(ResolveVersion)
	updRecord.CurrentVersion = newVersion
	condition := new(ResolveVersion)
	condition.Id = 1
	_, err := Engine.Table("resolve_config").Update(updRecord, condition)
	if err != nil {
		return false, err, 0, 0
	}
	return true, nil, oldVersion, newVersion
}

func RemoveResolveRecord(wrapper *ResolveRecord) (bool, error) {
	_, err := Engine.Table("resolve_record").Delete(wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func IsResolveRecordExist(wrapper *ResolveRecord) bool {
	wrapper.Version = GetResolveVersion()
	count, err := Engine.Table("resolve_record").Count(wrapper)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count > 0
}

func IsResolveRecordExistById(id int) bool {
	wrapper := new(ResolveRecord)
	wrapper.Id = id
	return IsResolveRecordExist(wrapper)
}

func IsUpdResolveRecordExist(id int, wrapper *ResolveRecord) bool {
	wrapper.Version = GetResolveVersion()
	count, err := Engine.Table("resolve_record").Where("id != ?", id).Count(wrapper)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count > 0
}

func ModifyResolveRecordById(id int, updateRecord *ResolveRecord) (bool, error) {
	updateRecord.UpdateTime = time.Now().Format(config.DataTimeFormat)
	wrapper := new(ResolveRecord)
	wrapper.Id = id
	_, err := Engine.Table("resolve_record").Update(updateRecord, wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func SwitchResolveRecord(id int, enabled int) (bool, error) {
	var updateRecord ResolveRecord
	updateRecord.UpdateTime = time.Now().Format(config.DataTimeFormat)
	updateRecord.Enabled = enabled
	wrapper := new(ResolveRecord)
	wrapper.Id = id
	_, err := Engine.Table("resolve_record").Update(updateRecord, wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func FindResolveVersion() []int {
	var records []int
	err := Engine.Select("distinct version").Cols("version").Table("resolve_record").Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}
