package code

type ResCode int64

// 这个文件专门放一些公共的响应码定义，或者一些工具函数，比如响应格式化等
// 定义一些响应码，注意：这些响应码是业务逻辑层的，不是 HTTP 状态码
// 你可以根据实际业务需求来定义这些响应码，下面是一些示例
// 1000-1999 是用户相关的响应码
const CtxUserIDKey = "user_id"

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserNotExist
	CodeUserAlreadyExists
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:           "success",
	CodeInvalidParam:      "请求参数错误",
	CodeUserAlreadyExists: "用户已存在",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "用户名或密码错误",
	CodeServerBusy:        "服务器繁忙",
	CodeInvalidToken:      "无效的Token",
	CodeNeedLogin:         "需要登录",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
