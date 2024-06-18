package utils

import "github.com/gin-gonic/gin"


type Resp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}


func ErrResponse(c *gin.Context,code int,massgae string,data any){
	c.JSON(200,Resp[any]{
		Code: code,
		Message: massgae,
		Data: data,
	})
}

func SuccessResp(c *gin.Context, data any) {
	c.JSON(200, Resp[any]{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}