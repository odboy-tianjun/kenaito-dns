package domain

/*
 * @Description  领域模型定义
 * @Author  www.odboy.cn
 * @Date  20241108
 */
type CreateResolveRecord struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Ttl   int    `json:"ttl" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type RemoveResolveRecord struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ModifyResolveRecord struct {
	Id    int    `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Ttl   int    `json:"ttl"`
	Value string `json:"value" binding:"required"`
}

type QueryPageArgs struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"pageSize" binding:"required"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type TestArgs struct {
	Name string `json:"name" binding:"required"`
}

type SwitchArgs struct {
	Id      int `json:"id" binding:"required"`
	Enabled int `json:"enabled" binding:"required"`
}

type QueryByIdArgs struct {
	Id int `json:"id" binding:"required"`
}

type RollbackVersionArgs struct {
	Version int `json:"version" binding:"required"`
}
