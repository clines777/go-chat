package protocol

import "github.com/gin-gonic/gin"

// ApiResponse http 返回 response
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewApiResponse(code int, msg string, data any) *ApiResponse {
	return &ApiResponse{Code: code, Message: msg, Data: data}
}

func (r *ApiResponse) H() map[string]any {
	return gin.H{"code": r.Code, "msg": r.Message, "data": r.Data}
}
