package cache

import (
	"fmt"
	"kenaito-dns/dao"
	"sync"
)

var KeyResolveRecordMap sync.Map
var IdResolveRecordMap sync.Map

func ReloadCache() {
	resolveRecords := dao.FindResolveRecordByVersion(dao.GetResolveVersion())
	for _, record := range resolveRecords {
		// id -> resolveRecord
		IdResolveRecordMap.Store(record.Id, record)
		// key -> resolveRecord
		cacheKey := fmt.Sprintf("%s-%s", record.Name, record.RecordType)
		records, ok := KeyResolveRecordMap.Load(cacheKey)
		if !ok {
			fmt.Println("读取缓存失败, key=" + cacheKey)
			var tempRecords []dao.ResolveRecord
			tempRecords = append(tempRecords, record)
			KeyResolveRecordMap.Store(cacheKey, tempRecords)
		} else {
			fmt.Println("读取缓存成功, key=" + cacheKey)
			var newRecords = records.([]dao.ResolveRecord)
			records = append(newRecords, record)
			KeyResolveRecordMap.Store(cacheKey, records)
		}
	}
}
