package service

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/lib"
)

type MemberService struct{}

func (userService *MemberService) GetMemberInfo(c *gin.Context) lib.ApiResponse {
	return lib.ApiResponse{}
}
