package logic

import (
	"BLUEBELL/db"
	"context"

	"go.uber.org/zap"
)

type CommunityLogic struct {
	store *db.Queries // 依赖注入
}

func NewCommunityLogic(store *db.Queries) *CommunityLogic {
	return &CommunityLogic{store: store}
}

func (l *CommunityLogic) GetCommunityList(ctx context.Context) ([]*db.GetCommunityListRow, error) {
	//查数据库，找到所有的社区，返回给 controller
	// 1. 查询数据库，获取社区列表
	communities, err := l.store.GetCommunityList(ctx) // 这里传入 context，如果需要的话
	// 2. 检查错误
	if err != nil {
		zap.L().Error("db.GetCommunityList() failed", zap.Error(err))
		return nil, err
	}
	// 3. 返回社区列表
	var result []*db.GetCommunityListRow
	for i := range communities {
		result = append(result, &communities[i]) // 注意：这里需要取地址，因为 GetCommunityList 返回的是值类型的切片
	}
	return result, nil
}

func (l *CommunityLogic) GetCommunityDetail(ctx context.Context, id int64) (*db.GetCommunityDetailByIDRow, error) {
	//查询数据库，获取社区详情
	data, err := l.store.GetCommunityDetailByID(ctx, uint64(id))
	if err != nil {
		zap.L().Error("db.GetCommunityDetailByID() failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	//返回社区详情
	return &data, nil
}
