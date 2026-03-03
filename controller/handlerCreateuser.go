package controller

import (
	"BLUEBELL/logic"
	"BLUEBELL/models"
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserHandler 定义一个结构体来持有 Logic 层对象
type UserHandler struct {
	userLogic *logic.UserLogic
}

// NewUserHandler 构造函数
func NewUserHandler(l *logic.UserLogic) *UserHandler {
	return &UserHandler{userLogic: l}
}

// SignUpHandler 现在作为 UserHandler 的方法
func (h *UserHandler) SignUpHandler(c *gin.Context) {
	// 1. 参数校验
	p := new(models.SignupParams)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUpHandler with invalid param", zap.Error(err))

		// 处理 validator 错误翻译... (保留你原本的逻辑)
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.RespondWithError(c, code.CodeInvalidParam) // 这里我们也返回 200，但 message 里有错误信息，保持和登录接口一致的风格
			return
		}
		response.RespondWithErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans))) // 这里我们也返回 200，但 message 里有错误信息，保持和登录接口一致的风格
		return
	}

	// 2. 业务处理 (调用 Logic 层)
	// 注意：这里使用的是 h.userLogic 调用方法
	if err := h.userLogic.SignUp(c.Request.Context(), p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, h.userLogic.CheckUserExists(c.Request.Context(), p.Username)) {
			response.RespondWithError(c, code.CodeUserAlreadyExists) // 业务错误也返回 200，但 message 里有错误信息
			return
		}
		response.RespondWithError(c, code.CodeServerBusy) // 业务错误也返回 200，但 message 里有错误信息
		return
	}

	// 3. 返回成功响应
	response.RespondWithSuccess(c, nil)
}
