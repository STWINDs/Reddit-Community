package controller

import (
	"BLUEBELL/models"
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func (l *PostHandler) PostVoteController(c *gin.Context) {
	p := new(models.VoteDataParam)
	if err := c.ShouldBindBodyWithJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.RespondWithError(c, code.CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		response.RespondWithErrorWithMsg(c, code.CodeInvalidParam, errData)
	}

	//获取UserID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("userID not found in context")
		response.RespondWithError(c, code.CodeNeedLogin)
		return
	}

	if err := l.logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost()Error")
	}
	response.RespondWithSuccess(c, nil)
}
