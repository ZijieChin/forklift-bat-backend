package util

import "github.com/gin-gonic/gin"

type Gin struct {
	Context *gin.Context
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var TimeInterval = 5 // seconds between two battery switch

func (g *Gin) Response(httpcode int, code int, msg string, data interface{}) {
	g.Context.JSON(httpcode, response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
