package entity

// Channel 通道账号配置
type Channel struct {
	Id             uint64
	Code           string //编码
	Type           string //类型。钉钉DINGDING；邮箱EMAIL；
	UserName       string //账号
	Password       string //密钥
	Address        string //请求地址
	Port           int    //端口
	CallbackUrl    string //回调地址
	QueryUrl       string //查询地址
	JsonConfig     string //json配置
	JsonConfigDesc string //json配置说明
	State          bool   //启用状态
	IsDefault      bool   //是否默认
}
