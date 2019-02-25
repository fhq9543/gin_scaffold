package controllers

import (
	"github.com/gin-gonic/gin"
	"go_base/utils/rescode"
	"net/http"
)

/**
全局成功\失败
**/
func APIResponse(c *gin.Context, success bool, data interface{}, msg string) {
	var resCode string
	if success {
		resCode = rescode.Success
		if msg == "" {
			msg = "操作成功"
		}
	} else {
		resCode = rescode.Error
		if msg == "" {
			msg = "操作失败"
		}
	}
	c.JSON(http.StatusOK, gin.H{"rescode": resCode, "data": data, "msg": msg})
}

//400
func APIResponseBadRequest(c *gin.Context, rescode string, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"rescode": rescode, "data": nil, "msg": msg})
}

//401
func APIResponseUnauthorized(c *gin.Context, rescode string, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{"rescode": rescode, "data": nil, "msg": msg})
}

//403
func APIResponseForbidden(c *gin.Context, rescode string, msg string) {
	c.JSON(http.StatusForbidden, gin.H{"rescode": rescode, "data": nil, "msg": msg})
}

//404
func APIResponseNotFound(c *gin.Context, rescode string, msg string) {
	c.JSON(http.StatusNotFound, gin.H{"rescode": rescode, "data": nil, "msg": msg})
}

//405
func APIResponseNotAllowed(c *gin.Context, rescode string, msg string) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"rescode": rescode, "data": nil, "msg": msg})
}

//406
func APIResponseNotAcceptable(c *gin.Context, rescode string, msg string) {
	c.JSON(http.StatusNotAcceptable, gin.H{"rescode": rescode, "data": nil, "msg": msg})
}

//500
func APIResponseException(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{"rescode": rescode.Error, "data": nil, "msg": msg})
}
