package controller

import (
	"BLUEBELL/models"
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/response"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func (h *UserHandler) LoginHandler(c *gin.Context) {
	// 1. 参数校验
	p := new(models.LoginParams)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("LoginHandler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.RespondWithError(c, code.CodeInvalidParam) // 这里我们也返回 200，但 message 里有错误信息，保持和注册接口一致的风格
			return
		}
		// 处理 validator 错误翻译... (保留你原本的逻辑)
		// 这里我们也返回 200，但 message 里有错误信息，保持和注册接口一致的风格
		response.RespondWithErrorWithMsg(c, code.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	// 2. 业务处理 (调用 Logic 层)
	// 注意：这里使用的是 h.userLogic 调用方法
	user, err := h.userLogic.Login(c.Request.Context(), p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, h.userLogic.CheckUserExists(c.Request.Context(), p.Username)) {
			response.RespondWithError(c, code.CodeUserNotExist) // 业务错误也返回 200，但 message 里有错误信息
			return
		}
		response.RespondWithError(c, code.CodeInvalidPassword) // 业务错误也返回 200，但 message 里有错误信息
		return
	}

	// 3. 返回成功响应，包含 token
	response.RespondWithSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
