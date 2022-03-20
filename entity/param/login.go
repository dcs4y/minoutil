package param

type LoginRequest struct {
	Type      string `json:"type"`      // 登录类型：1(后台管理员)，2(前端用户)。
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

func (t LoginRequest) Verify() error {
	return nil
}
