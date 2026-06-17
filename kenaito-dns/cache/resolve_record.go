package cache

/*
 * @Description  缓存
 * @Author  https://www.odboy.cn
 * @Date  20241107
 */
import (
	"fmt"
	"kenaito-dns/common"
	"kenaito-dns/dao"
	"kenaito-dns/util"
	"sync"
)

var KeyResolveRecordMap sync.Map
var IdResolveRecordMap sync.Map
var NameSet *common.Set

var reloadMu sync.Mutex

func ReloadCache() {
	reloadMu.Lock()
	defer reloadMu.Unlock()

	fmt.Println("[app]  [info]  " + util.NowStr() + " [Cache] Reload cache start")
	// 清空旧缓存
	KeyResolveRecordMap.Range(func(key, value any) bool {
		KeyResolveRecordMap.Delete(key)
		return true
	})
	IdResolveRecordMap.Range(func(key, value any) bool {
		IdResolveRecordMap.Delete(key)
		return true
	})
	NameSet = common.NewSet()
	resolveRecords := dao.FindResolveRecordByVersion(dao.GetResolveVersion(), false)
	for _, record := range resolveRecords {
		IdResolveRecordMap.Store(record.Id, record)
		cacheKey := fmt.Sprintf("%s-%s", record.Name, record.RecordType)
		records, ok := KeyResolveRecordMap.Load(cacheKey)
		if ok {
			recordSlice := records.([]dao.ResolveRecord)
			recordSlice = append(recordSlice, record)
			KeyResolveRecordMap.Store(cacheKey, recordSlice)
		} else {
			var tempRecords []dao.ResolveRecord
			tempRecords = append(tempRecords, record)
			KeyResolveRecordMap.Store(cacheKey, tempRecords)
		}
		NameSet.Add(record.Name)
	}
	fmt.Println("[app]  [info]  " + util.NowStr() + " [Cache] Reload cache end")
}
