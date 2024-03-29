package baseModel

// PageParam 分页参数
type PageParam struct {
	PageStart  int `gorm:"-"`
	PageSize   int `json:"pageSize" form:"pageSize" gorm:"-"`
	PageNumber int `json:"pageNumber" form:"pageNumber" gorm:"-"`
}

func (pg PageParam) GetStart() int {
	start := 0
	if pg.PageNumber > 0 {
		start = (pg.PageNumber - 1) * pg.PageSize
	}
	return start
}

// IdParam 通用ID参数
type IdParam struct {
	Id  int64   `json:"id"`
	Ids []int64 `json:"ids"`
}
