package result

type CaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength""`
}

type LoginResponse struct {
	Username    string `json:"userName"`    // 用户登录名
	Password    string `json:"-"`           // 用户登录密码
	NickName    string `json:"nickName"`    // 用户昵称
	SideMode    string `json:"sideMode"`    // 用户侧边主题
	HeaderImg   string `json:"headerImg"`   // 用户头像
	BaseColor   string `json:"baseColor"`   // 基础颜色
	ActiveColor string `json:"activeColor"` // 活跃颜色
	//Authority   SysAuthority   `json:"authority"` // 用户角色
	Phone     string `json:"phone"` // 用户手机号
	Email     string `json:"email"` // 用户邮箱
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}
