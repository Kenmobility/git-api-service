package response

import "github.com/gin-gonic/gin"

type ResponseObject struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success is used to send a successful request response
func Success(c *gin.Context, code int, message string, data interface{}) {
	obj := ResponseObject{
		Code:    code,
		Message: message,
	}
	if data != nil {
		obj.Data = data
	}

	c.JSON(code, obj)
}

// Failure is used to send a failed request response
func Failure(c *gin.Context, code int, message string, error interface{}) {
	obj := ResponseObject{
		Code:    code,
		Message: message,
		Error:   error,
	}
	c.JSON(code, obj)
}
