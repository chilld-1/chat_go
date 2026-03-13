package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r Response) BadRequestResponse(c *gin.Context, s string) {
	panic("unimplemented")
}

func NewResponse(c *gin.Context, code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewResponse(c, 0, "success", data))
}
func ErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, NewResponse(c, 1, msg, nil))
}
func BadRequestResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, NewResponse(c, http.StatusBadRequest, msg, nil))
}
func UnauthorizedResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, NewResponse(c, http.StatusUnauthorized, msg, nil))
}
func ForbiddenResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, NewResponse(c, http.StatusForbidden, msg, nil))
}
func NotFoundResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, NewResponse(c, http.StatusNotFound, msg, nil))
}
func InternalServerErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, NewResponse(c, http.StatusInternalServerError, msg, nil))
}
