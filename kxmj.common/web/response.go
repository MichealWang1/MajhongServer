package web

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RespSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: codes.Success,
		Msg:  codes.GetMessage(codes.Success),
		Data: data,
	})
}

func RespFailed(ctx *gin.Context, code int, msg ...string) {
	var m string
	if len(msg) <= 0 {
		m = codes.GetMessage(code)
	} else {
		m = msg[0]
	}
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  m,
		Data: nil,
	})
}

func RespFailWithErr(ctx *gin.Context, err error) {
	e := codes.ParseError(err)
	RespFailed(ctx, e.Code, e.Message)
}
