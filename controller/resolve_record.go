package controller

/*
 * @Description  Web控制层
 * @Author  www.odboy.cn
 * @Date  20241108
 */
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kenaito-dns/cache"
	"kenaito-dns/constant"
	"kenaito-dns/dao"
	"kenaito-dns/domain"
	"kenaito-dns/util"
	"net/http"
	"strings"
)

func InitRestFunc(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    "pong",
		})
	})
	// 测试解析状态
	r.POST("/test", func(c *gin.Context) {
		var jsonObj domain.TestArgs
		err := c.ShouldBindJSON(&jsonObj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
			return
		}
		name := jsonObj.Name
		valid := util.IsValidDomain(name)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"message": "域名解析失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "域名解析成功",
		})
	})
	// 启停记录
	r.POST("/switch", func(c *gin.Context) {
		var jsonObj domain.SwitchArgs
		err := c.ShouldBindJSON(&jsonObj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
			return
		}
		_, err = dao.SwitchResolveRecord(jsonObj.Id, jsonObj.Enabled)
		if err != nil {
			if jsonObj.Enabled == 1 {
				c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("启用失败, %v", err)})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("停用失败, %v", err)})
			}
			return
		}
		cache.ReloadCache()
		if jsonObj.Enabled == 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "启用成功",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "停用成功",
			})
		}
	})
	// 创建RR记录
	r.POST("/create", func(c *gin.Context) {
		var jsonObj domain.CreateResolveRecord
		err := c.ShouldBindJSON(&jsonObj)
		newRecord, isErr := validRequestBody(c, err, jsonObj.Name, jsonObj.Type, jsonObj.Ttl, jsonObj.Value, false)
		if isErr {
			return
		}
		if dao.IsResolveRecordExist(newRecord) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "记录 " + newRecord.Name + " " + newRecord.RecordType + " " + newRecord.Value + " 已存在",
			})
			return
		}
		newRecord.Ttl = jsonObj.Ttl
		executeResult, err, oldVersion, newVersion := dao.BackupResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("备份"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		executeResult, _ = dao.SaveResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("添加"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		cache.ReloadCache()
		body := make(map[string]interface{})
		body["oldVersion"] = oldVersion
		body["newVersion"] = newVersion
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "添加" + newRecord.RecordType + "记录成功",
			"data":    body,
		})
		return
	})
	// 删除RR记录
	r.POST("/remove", func(c *gin.Context) {
		var jsonObj domain.RemoveResolveRecord
		err := c.ShouldBindJSON(&jsonObj)
		newRecord, isErr := validRequestBody(c, err, jsonObj.Name, jsonObj.Type, 0, jsonObj.Value, true)
		if isErr {
			return
		}
		if !dao.IsResolveRecordExist(newRecord) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "记录 " + newRecord.Name + " " + newRecord.RecordType + " " + newRecord.Value + " 不存在"})
			return
		}
		executeResult, err, oldVersion, newVersion := dao.BackupResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("备份"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		executeResult, err = dao.RemoveResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("删除"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		cache.ReloadCache()
		body := make(map[string]interface{})
		body["oldVersion"] = oldVersion
		body["newVersion"] = newVersion
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "删除" + newRecord.RecordType + "记录成功",
			"data":    body,
		})
		return
	})
	// 修改RR记录
	r.POST("/modify", func(c *gin.Context) {
		var jsonObj domain.ModifyResolveRecord
		err := c.ShouldBindJSON(&jsonObj)
		newRecord, isErr := validModifyRequestBody(c, err, jsonObj.Name, jsonObj.Type, jsonObj.Ttl, jsonObj.Value)
		if isErr {
			return
		}
		if jsonObj.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数ID(id)必填"})
			return
		}
		if !dao.IsResolveRecordExistById(jsonObj.Id) {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("记录 id=%d 不存在", jsonObj.Id)})
			return
		}
		if dao.IsUpdResolveRecordExist(jsonObj.Id, newRecord) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "记录 " + jsonObj.Name + " " + jsonObj.Type + " " + jsonObj.Value + " 已存在"})
			return
		}
		localOldRecord := dao.FindResolveRecordById(jsonObj.Id)
		executeResult, err, oldVersion, newVersion := dao.BackupResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("备份"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		localNewRecord := dao.FindOneResolveRecord(localOldRecord, newVersion)
		if localNewRecord == nil || localNewRecord.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("查询待更新" + newRecord.RecordType + "记录失败")})
			return
		}
		if localNewRecord.Ttl == 0 {
			localNewRecord.Ttl = 10
		}
		updRecord := new(dao.ResolveRecord)
		updRecord.Name = newRecord.Name
		updRecord.RecordType = newRecord.RecordType
		updRecord.Ttl = newRecord.Ttl
		updRecord.Value = newRecord.Value
		updRecord.CreateTime = localNewRecord.CreateTime
		executeResult, err = dao.ModifyResolveRecordById(localNewRecord.Id, updRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("更新"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		cache.ReloadCache()
		body := make(map[string]interface{})
		body["oldVersion"] = oldVersion
		body["newVersion"] = newVersion
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "更新" + newRecord.RecordType + "记录成功",
			"data":    body,
		})
		return
	})
	// 分页查询RR记录
	r.POST("/queryPage", func(c *gin.Context) {
		var jsonObj domain.QueryPageArgs
		err := c.ShouldBindJSON(&jsonObj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
			return
		}
		records := dao.FindResolveRecordPage(jsonObj.Page, jsonObj.PageSize, &jsonObj)
		count := dao.CountResolveRecordPage(jsonObj.Page, jsonObj.PageSize, &jsonObj)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "分页查询RR记录成功",
			"data":    records,
			"count":   count,
		})
		return
	})
	// 根据id查询RR记录明细
	r.POST("/queryById", func(c *gin.Context) {
		var jsonObj domain.QueryByIdArgs
		err := c.ShouldBindJSON(&jsonObj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
			return
		}
		records := dao.FindResolveRecordById(jsonObj.Id)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "根据id查询RR记录明细成功",
			"data":    records,
		})
		return
	})
	// 分页查询变更历史记录
	r.POST("/queryVersionPage", func(c *gin.Context) {
		var jsonObj domain.QueryPageArgs
		err := c.ShouldBindJSON(&jsonObj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
			return
		}
		records := dao.FindResolveVersionPage(jsonObj.Page, jsonObj.PageSize)
		count := dao.CountResolveVersionPage(jsonObj.Page, jsonObj.PageSize)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "分页查询变更历史记录成功",
			"data":    records,
			"count":   count,
		})
		return
	})
	// 回滚到某一版本
	r.POST("/rollback", func(c *gin.Context) {
		var jsonObj domain.RollbackVersionArgs
		err := c.ShouldBindJSON(&jsonObj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
			return
		}
		versions := dao.FindResolveRecordByVersion(jsonObj.Version, true)
		if len(versions) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("版本号 %d 不存在, 回滚失败", jsonObj.Version)})
			return
		}
		//executeResult, err := dao.ModifyResolveVersion(jsonObj.Version)
		//if !executeResult {
		//	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("回滚失败, %v", err)})
		//	return
		//}
		//cache.ReloadCache()
		//body := make(map[string]interface{})
		//body["currentVersion"] = jsonObj.Version
		//c.JSON(http.StatusOK, gin.H{
		//	"code":    0,
		//	"message": "回滚成功",
		//	"data":    body,
		//})
		return
	})
}
func validRequestBody(c *gin.Context, err error, name string, recordType string, ttl int, value string, isDelete bool) (*dao.ResolveRecord, bool) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
		return nil, true
	}
	if util.IsBlank(name) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数主机记录(name)必填, 例如: www.odboy.cn 或 odboy.cn"})
		return nil, true
	}
	if util.IsBlank(recordType) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录类型(type)必填, 目前支持 A、AAAA 记录"})
		return nil, true
	}
	if util.IsBlank(value) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)必填"})
		return nil, true
	}
	if !isDelete {
		if ttl < 10 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数缓存有效时间(ttl)有误，必须大于等于10"})
			return nil, true
		}
	}
	if !util.IsValidName(name) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数主机记录(name)有误，无效的主机记录"})
		return nil, true
	}
	switch recordType {
	case constant.R_A:
		if !util.IsIPv4(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的IPv4地址"})
			return nil, true
		}
	case constant.R_AAAA:
		if !util.IsIPv6(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的IPv6地址"})
			return nil, true
		}
	case constant.R_CNAME:
		if !util.IsValidDomain(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的主机记录"})
			return nil, true
		}
	case constant.R_MX:
		if !util.IsValidDomain(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的主机记录"})
			return nil, true
		}
	case constant.R_TXT:
		if len(value) > 512 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，长度必须 <= 512"})
			return nil, true
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(type)有误，不支持的记录类型: " + recordType})
		return nil, true
	}
	newRecord := new(dao.ResolveRecord)
	newRecord.Name = strings.TrimSpace(name)
	newRecord.RecordType = strings.TrimSpace(recordType)
	newRecord.Value = strings.TrimSpace(value)
	return newRecord, false
}
func validModifyRequestBody(c *gin.Context, err error, name string, recordType string, ttl int, value string) (*dao.ResolveRecord, bool) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("校验失败, %v", err)})
		return nil, true
	}
	if ttl != 0 && ttl < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数缓存有效时间(ttl)有误，必须大于等于10"})
		return nil, true
	}
	if !util.IsValidName(name) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数主机记录(name)有误，无效的主机记录"})
		return nil, true
	}
	switch recordType {
	case constant.R_A:
		if !util.IsIPv4(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的IPv4地址"})
			return nil, true
		}
	case constant.R_AAAA:
		if !util.IsIPv6(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的IPv6地址"})
			return nil, true
		}
	case constant.R_CNAME:
		if !util.IsValidDomain(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的主机记录"})
			return nil, true
		}
	case constant.R_MX:
		if !util.IsValidDomain(value) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，无效的主机记录"})
			return nil, true
		}
	case constant.R_TXT:
		if len(value) > 512 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(value)有误，长度必须 <= 512"})
			return nil, true
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数记录值(type)有误，不支持的记录类型: " + recordType})
		return nil, true
	}
	newRecord := new(dao.ResolveRecord)
	newRecord.Name = strings.TrimSpace(name)
	newRecord.RecordType = strings.TrimSpace(recordType)
	newRecord.Value = strings.TrimSpace(value)
	return newRecord, false
}
