package baseModel

import "errors"

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
	Id  string   `json:"id"`
	Ids []string `json:"ids"`
}

func (param IdParam) Verify() error {
	if param.Id == "" {
		if len(param.Ids) == 0 {
			return errors.New("ID不能为空！")
		} else {
			for _, id := range param.Ids {
				if id == "" {
					return errors.New("ID不能为空！")
				}
			}
		}
	}
	return nil
}
