package handler

import (
	"game/entity/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取登录菜单
func getMenu(c *gin.Context) {
	var menus []model.Menu
	c.JSON(http.StatusOK, menus)
}
