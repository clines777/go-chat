package controller

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/lib"
	"gochat/internal/service"
)

type LoginController struct {
	service *service.LoginService
}

func (controller *LoginController) Login(c *gin.Context) lib.ApiResponse {
	return lib.ApiResponse{}
}
