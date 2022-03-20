package common

const (
	CodeSuccess = 0
	CodeError   = 10001
)

// Result 统计返回格式
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// ResultOkData 返回成功数据
func ResultOkData(date interface{}) (result Result) {
	result.Code = CodeSuccess
	result.Message = "OK"
	result.Data = date
	return
}

// ResultOkMessage 返回成功消息
func ResultOkMessage(message string) (result Result) {
	result.Code = CodeSuccess
	result.Message = message
	return
}

// ResultOk 返回成功的消息及数据
func ResultOk(message string, date interface{}) (result Result) {
	result.Code = CodeSuccess
	result.Message = message
	result.Data = date
	return
}

// ResultErrorMessage 返回错误信息
func ResultErrorMessage(message string) (result Result) {
	result.Code = CodeError
	result.Message = message
	return
}

// ResultError 错误码及消息
func ResultError(code int, message string) (result Result) {
	result.Code = code
	result.Message = message
	return
}
