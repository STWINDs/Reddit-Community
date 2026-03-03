package controller

import (
	"BLUEBELL/logic"
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommunityHandler struct {
	logic *logic.CommunityLogic // 持有对应的社区逻辑层
}

func NewCommunityHandler(l *logic.CommunityLogic) *CommunityHandler {
	return &CommunityHandler{logic: l}
}

// --- 跟社区相关的---
func (h *CommunityHandler) GetCommunityHandler(c *gin.Context) {
	// 1. 从数据库获取社区列表，查询到所有的社区：community_id, community_name，以列表的形式返回
	communities, err := h.logic.GetCommunityList(c.Request.Context()) // 这里传入 context，如果需要的话

	// 2. 检查错误
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy) //不轻易吧服务器错误暴露给用户，统一返回一个服务器忙的错误
		return
	}

	// 3. 返回社区列表
	response.RespondWithSuccess(c, communities)
}

// 根据ID查询社区分类详情
func (h *CommunityHandler) GetCommunityDetailHandler(c *gin.Context) {
	//获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("invalid community id", zap.String("id", idStr), zap.Error(err))
		response.RespondWithError(c, code.CodeInvalidParam)
		return
	}
	//查询数据库，获取社区详情
	data, err := h.logic.GetCommunityDetail(c.Request.Context(), id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Int64("id", id), zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy)
		return
	}
	//检查错误
	//返回社区详情
	response.RespondWithSuccess(c, data)

}
