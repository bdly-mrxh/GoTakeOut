package vo

// PageResult 分页结果
type PageResult struct {
	Total   int64 `json:"total"`   // 总记录数
	Records any   `json:"records"` // 当前页数据列表
}
