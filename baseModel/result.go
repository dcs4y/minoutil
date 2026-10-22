package baseModel

const (
	CodeSuccess = 0
	TokenError  = 1000
	CodeError   = 9999
)

// Result 统计返回格式
type Result[D any] struct {
	Code    int    `json:"code"` // 状态码
	Message string `json:"msg"`  // 返回信息
	Data    D      `json:"data"` // 返回数据
}

type PageResult[D any, S any] struct {
	Code    int    `json:"code"` // 状态码
	Message string `json:"msg"`  // 返回信息
	Data    struct {
		List    []D   `json:"list"`    //列表数据
		Total   int64 `json:"total"`   //总数
		Summary S     `json:"summary"` //统计数据
	} `json:"data"` // 返回数据
}

func (result Result[D]) Ok(data D) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data = data
}

func (result Result[D]) Error(message string) {
	result.Code = CodeError
	result.Message = message
}

func (result PageResult[D, S]) Ok(data []D, total int64) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data.List = data
	result.Data.Total = total
}

func (result PageResult[D, S]) Error(message string) {
	result.Code = CodeError
	result.Message = message
}

// ResultOkData 返回成功数据
func ResultOkData[D any](data D) (result Result[D]) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data = data
	return
}

// ResultOkPageData 返回成功的分页数据
/*func ResultOkPageData(list any, total int64) (result Result) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data = map[string]any{
		"list":  list,
		"total": total,
	}
	return
}*/

func ResultOkPageData[D any, S any](list []D, total int64, summary S) (result PageResult[D, S]) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data.List = list
	result.Data.Total = total
	result.Data.Summary = summary
	return
}

func ResultErrorPage(message string) (result PageResult[any, any]) {
	result.Code = CodeError
	result.Message = message
	return
}

// ResultOk 返回默认成功信息
func ResultOk() (result Result[any]) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	return
}

// ResultOkMessage 返回成功消息
func ResultOkMessage(message string) (result Result[any]) {
	result.Code = CodeSuccess
	result.Message = message
	return
}

// ResultOkDataMessage 返回成功的数据及消息
func ResultOkDataMessage[D any](data D, message string) (result Result[D]) {
	result.Code = CodeSuccess
	result.Message = message
	result.Data = data
	return
}

// ResultErrorMessage 返回错误信息
func ResultErrorMessage(message string) (result Result[any]) {
	result.Code = CodeError
	result.Message = message
	return
}

// ResultErrorReload 返回错误信息并退出登录状态
func ResultErrorReload[D map[string]bool](message string) (result Result[D]) {
	result.Code = TokenError
	result.Message = message
	result.Data = map[string]bool{
		"reload": true,
	}
	return
}

// ResultError 错误码及消息
func ResultError(code int, message string) (result Result[any]) {
	result.Code = code
	result.Message = message
	return
}
