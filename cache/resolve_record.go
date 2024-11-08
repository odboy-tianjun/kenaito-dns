package cache

import (
	"fmt"
	"kenaito-dns/config"
	"kenaito-dns/dao"
	"sync"
	"time"
)

var KeyResolveRecordMap sync.Map
var IdResolveRecordMap sync.Map

func ReloadCache() {
	fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Reload cache start")
	resolveRecords := dao.FindResolveRecordByVersion(dao.GetResolveVersion())
	for _, record := range resolveRecords {
		// id -> resolveRecord
		IdResolveRecordMap.Store(record.Id, record)
		// key -> resolveRecord
		cacheKey := fmt.Sprintf("%s-%s", record.Name, record.RecordType)
		records, ok := KeyResolveRecordMap.Load(cacheKey)
		if !ok {
			var tempRecords []dao.ResolveRecord
			tempRecords = append(tempRecords, record)
			KeyResolveRecordMap.Store(cacheKey, tempRecords)
		} else {
			var newRecords = records.([]dao.ResolveRecord)
			records = append(newRecords, record)
			KeyResolveRecordMap.Store(cacheKey, records)
		}
	}
	fmt.Println("[app]  [info]  " + time.Now().Format(config.AppTimeFormat) + " [Cache] Reload cache end")
}
