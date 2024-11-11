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
	KeyResolveRecordMap.Range(cleanKeyCache)
	IdResolveRecordMap.Range(cleanIdCache)
	resolveRecords := dao.FindResolveRecordByVersion(dao.GetResolveVersion(), false)
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

func cleanKeyCache(key any, value any) bool {
	KeyResolveRecordMap.Delete(key)
	return true
}

func cleanIdCache(key any, value any) bool {
	IdResolveRecordMap.Delete(key)
	return true
}
