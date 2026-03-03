package response

import (
	"BLUEBELL/pkg/code"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code code.ResCode `json:"code"`
	Msg  interface{}  `json:"msg"`
	Data interface{}  `json:"data,omitempty"`
}

// 这个文件专门放一些公共的响应函数，或者一些工具函数，比如响应格式化等
func RespondWithError(c *gin.Context, code code.ResCode) {
	// 这里我们统一返回 200，但 message 里有错误信息，保持和注册接口一致的风格
	// 为什么要用&ResponseData 而不是直接 ResponseData？因为我们需要传递一个指针给 c.JSON，这样 Gin 才能正确地序列化它。
	// 如果直接传递 ResponseData 的值，Gin 可能会遇到一些问题，因为它期望一个指针类型来进行序列化。
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func RespondWithErrorWithMsg(c *gin.Context, code code.ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code.CodeSuccess,
		Msg:  code.CodeSuccess.Msg(),
		Data: data,
	})
}
