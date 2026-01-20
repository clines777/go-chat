package lib

import "github.com/gin-gonic/gin"

// ApiResponse http 返回 response
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *ApiResponse) Get() map[string]any {
	return gin.H{"code": r.Code, "message": r.Message, "data": r.Data}
}
