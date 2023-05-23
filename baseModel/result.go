package baseModel

const (
	CodeSuccess = 0
	CodeError   = 9999
)

// Result 统计返回格式
type Result struct {
	Code    int         `json:"code"` // 状态码
	Message string      `json:"msg"`  // 返回信息
	Data    interface{} `json:"data"` // 返回数据
}

// ResultOkData 返回成功数据
func ResultOkData(data interface{}) (result Result) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data = data
	return
}

// ResultOkPageData 返回成功的分页数据
func ResultOkPageData(list interface{}, total int64) (result Result) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	result.Data = map[string]interface{}{
		"list":  list,
		"total": total,
	}
	return
}

// ResultOk 返回默认成功信息
func ResultOk() (result Result) {
	result.Code = CodeSuccess
	result.Message = "操作成功!"
	return
}

// ResultOkMessage 返回成功消息
func ResultOkMessage(message string) (result Result) {
	result.Code = CodeSuccess
	result.Message = message
	return
}

// ResultOkDataMessage 返回成功的数据及消息
func ResultOkDataMessage(data interface{}, message string) (result Result) {
	result.Code = CodeSuccess
	result.Message = message
	result.Data = data
	return
}

// ResultErrorMessage 返回错误信息
func ResultErrorMessage(message string) (result Result) {
	result.Code = CodeError
	result.Message = message
	return
}

// ResultErrorReload 返回错误信息并退出登录状态
func ResultErrorReload(message string) (result Result) {
	result.Code = CodeError
	result.Message = message
	result.Data = map[string]bool{
		"reload": true,
	}
	return
}

// ResultError 错误码及消息
func ResultError(code int, message string) (result Result) {
	result.Code = code
	result.Message = message
	return
}
