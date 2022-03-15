package common

// PageParam 分页参数
type PageParam struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}
