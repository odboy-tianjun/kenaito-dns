package controller

/*
 * @Description  Web控制层
 * @Author  www.odboy.cn
 * @Date  20241108
 */
import (
	"fmt"
	"github.com/gin-gonic/gin"
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
			"message": "pong",
		})
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "记录 " + newRecord.Name + " " + newRecord.RecordType + " " + newRecord.Value + " 已存在"})
			return
		}
		newRecord.Ttl = jsonObj.Ttl
		executeResult, err, oldVersion, newVersion := dao.BackupResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("添加"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		executeResult, _ = dao.SaveResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("添加"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		body := make(map[string]interface{})
		body["oldVersion"] = oldVersion
		body["newVersion"] = newVersion
		c.JSON(http.StatusOK, gin.H{
			"message": "添加" + newRecord.RecordType + "记录成功",
			"body":    body,
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
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("删除"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		executeResult, err = dao.RemoveResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("删除"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		body := make(map[string]interface{})
		body["oldVersion"] = oldVersion
		body["newVersion"] = newVersion
		c.JSON(http.StatusOK, gin.H{
			"message": "删除" + newRecord.RecordType + "记录成功",
			"body":    body,
		})
		return
	})
	// 修改RR记录
	r.POST("/modify", func(c *gin.Context) {
		var jsonObj domain.ModifyResolveRecord
		err := c.ShouldBindJSON(&jsonObj)
		if jsonObj.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "参数ID(id)必填"})
			return
		}
		newRecord, isErr := validRequestBody(c, err, jsonObj.Name, jsonObj.Type, jsonObj.Ttl, jsonObj.Value, false)
		if isErr {
			return
		}
		if !dao.IsResolveRecordExistById(jsonObj.Id) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "记录 " + newRecord.Name + " " + newRecord.RecordType + " " + newRecord.Value + " 不存在"})
			return
		}
		if dao.IsUpdResolveRecordExist(jsonObj.Id, newRecord) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "记录 " + newRecord.Name + " " + newRecord.RecordType + " " + newRecord.Value + " 已存在"})
			return
		}
		executeResult, err, oldVersion, newVersion := dao.BackupResolveRecord(newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("修改"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		executeResult, err = dao.ModifyResolveRecordById(jsonObj.Id, newRecord)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("更新"+newRecord.RecordType+"记录失败, %v", err)})
			return
		}
		body := make(map[string]interface{})
		body["oldVersion"] = oldVersion
		body["newVersion"] = newVersion
		c.JSON(http.StatusOK, gin.H{
			"message": "更新" + newRecord.RecordType + "记录成功",
			"body":    body,
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
		c.JSON(http.StatusOK, gin.H{"message": "分页查询RR记录成功", "body": records})
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
		c.JSON(http.StatusOK, gin.H{"message": "根据id查询RR记录明细成功", "body": records})
		return
	})
	// 查询变更历史记录
	r.POST("/queryVersionList", func(c *gin.Context) {
		records := dao.FindResolveVersion()
		c.JSON(http.StatusOK, gin.H{"message": "查询变更历史记录列表成功", "body": records})
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
		versions := dao.FindResolveRecordByVersion(jsonObj.Version)
		if len(versions) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("版本号 %d 不存在, 回滚失败", jsonObj.Version)})
			return
		}
		executeResult, err := dao.ModifyResolveVersion(jsonObj.Version)
		if !executeResult {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("回滚失败, %v", err)})
			return
		}
		body := make(map[string]interface{})
		body["currentVersion"] = jsonObj.Version
		c.JSON(http.StatusOK, gin.H{
			"message": "回滚成功",
			"body":    body,
		})
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
