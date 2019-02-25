package request

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequestValidate(c *gin.Context, valid validation.Validation) {
	errors := gin.H{}
	for _, err := range valid.Errors {
		fmt.Println(err.Field, err.Message)
		errors[strings.ToLower(err.Field)] = err.Message
	}
	c.JSON(http.StatusBadRequest, errors)
}
