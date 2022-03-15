package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
	"net/http"
	"time"
)

//https://github.com/storyicon/grbac/blob/master/docs/README-chinese.md
//Grbac是一个快速，优雅和简洁的RBAC框架。它支持增强的通配符并使用Radix树匹配HTTP请求。令人惊奇的是，您可以在任何现有的数据库和数据结构中轻松使用它。
//grbac的作用是确保指定的资源只能由指定的角色访问。请注意，grbac不负责存储鉴权规则和分辨“当前请求发起者具有哪些角色”，更不负责角色的创建、分配等。这意味着您应该首先配置规则信息，并提供每个请求的发起者具有的角色。
//grbac将Host、Path和Method的组合视为Resource，并将Resource绑定到一组角色规则（称为Permission）。只有符合这些规则的用户才能访问相应的Resource。
//读取鉴权规则的组件称为Loader。grbac预置了一些Loader，你也可以通过实现func()(grbac.Rules，error)来根据你的设计来自定义Loader，并通过grbac.WithLoader加载它。

func LoadAuthorizationRules() (rules grbac.Rules, err error) {
	// 在这里实现你的逻辑
	// ...
	// 你可以从数据库或文件加载授权规则
	// 但是你需要以 grbac.Rules 的格式返回你的身份验证规则
	// 提示：你还可以将此函数绑定到golang结构体
	return
}

func QueryRolesByHeaders(header http.Header) (roles []string, err error) {
	// 在这里实现你的逻辑
	// ...
	// 这个逻辑可能是从请求的Headers中获取token，并且根据token从数据库中查询用户的相应角色。
	return roles, err
}

func RBACAuth() gin.HandlerFunc {
	// 在这里，我们通过“grbac.WithLoader”接口使用自定义Loader功能
	// 并指定应每分钟调用一次LoadAuthorizationRules函数以获取最新的身份验证规则。
	// Grbac还提供一些现成的Loader：
	// grbac.WithYAML
	// grbac.WithRules
	// grbac.WithJSON
	// ...
	rbac, err := grbac.New(grbac.WithLoader(LoadAuthorizationRules, time.Minute))
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		roles, err := QueryRolesByHeaders(c.Request.Header)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		state, _ := rbac.IsRequestGranted(c.Request, roles)
		if !state.IsGranted() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
