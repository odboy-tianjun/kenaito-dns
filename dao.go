package main

import "fmt"

type ResolveRecord struct {
	Id         int32  `xorm:"pk not null integer 'id' autoincr"`
	Name       string `xorm:"not null text 'name'"`
	RecordType string `xorm:"not null text 'record_type'"`
	Ttl        int32  `xorm:"not null integer 'ttl'"`
	Value      string `xorm:"not null text 'value'"`
}

func selectResolveRecords() []*ResolveRecord {
	records := make([]*ResolveRecord, 0)
	err := Engine.Table("resolve_record").Find(&records)
	if err != nil {
		fmt.Println(err)
	}
	return records
}
