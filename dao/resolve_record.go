package dao

/*
 * @Description  解析记录定义与操作
 * @Author  www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"kenaito-dns/domain"
	"kenaito-dns/util"
	"strings"
)

type ResolveRecord struct {
	Id         int    `xorm:"pk not null integer 'id' autoincr"`
	Name       string `xorm:"not null text 'name'"`
	RecordType string `xorm:"not null text 'record_type'"`
	Ttl        int    `xorm:"not null integer 'ttl'"`
	Value      string `xorm:"not null text 'value'"`
	Version    int    `xorm:"not null integer 'version'"`
}

func FindResolveRecordById(id int) []ResolveRecord {
	var records []ResolveRecord
	err := Engine.Table("resolve_record").Where("`id` = ?", id).Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}

func FindResolveRecordByVersion(version int) []ResolveRecord {
	var records []ResolveRecord
	err := Engine.Table("resolve_record").Where("`version` = ?", version).Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}

func FindResolveRecordByNameType(name string, recordType string) []ResolveRecord {
	var records []ResolveRecord
	err := Engine.Table("resolve_record").Where("`name` = ? and `record_type` = ? and `version` = ?", name, recordType, getResolveVersion()).Find(&records)
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
	session.And("`version` = ?", getResolveVersion())
	err := session.Limit(pageSize, offset).Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}

func SaveResolveRecord(wrapper *ResolveRecord) (bool, error) {
	_, err := Engine.Table("resolve_record").Insert(wrapper)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func BackupResolveRecord(record *ResolveRecord) (bool, error, int, int) {
	var backupRecords []*ResolveRecord
	oldVersion := getResolveVersion()
	newVersion := getResolveVersion() + 1
	oldRecords := FindResolveRecordByVersion(oldVersion)
	for _, oldRecord := range oldRecords {
		newRecord := new(ResolveRecord)
		newRecord.Name = oldRecord.Name
		newRecord.RecordType = oldRecord.RecordType
		newRecord.Ttl = oldRecord.Ttl
		newRecord.Value = oldRecord.Value
		newRecord.Version = newVersion
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
	wrapper.Version = getResolveVersion()
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
	r := new(ResolveRecord)
	r.Name = wrapper.Name
	r.RecordType = wrapper.RecordType
	r.Value = wrapper.Value
	r.Version = getResolveVersion()
	count, err := Engine.Table("resolve_record").Where("id != ?", id).Count(r)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count > 0
}

func ModifyResolveRecordById(id int, updateRecord *ResolveRecord) (bool, error) {
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
