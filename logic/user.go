package logic

import (
	"BLUEBELL/db"
	"BLUEBELL/models"
	"BLUEBELL/pkg/jwt"
	"BLUEBELL/pkg/snowflake"

	"context"
	"database/sql"
	"errors"
)

type UserLogic struct {
	store *db.Queries // 依赖注入
}

func NewUserLogic(store *db.Queries) *UserLogic {
	return &UserLogic{store: store}
}

// CheckUserExists 检查用户是否存在
func (l *UserLogic) CheckUserExists(ctx context.Context, username string) error {
	// 调用 sqlc 生成的 GetUserByUsername
	// 我们只关心 error，不关心返回的用户数据
	_, err := l.store.GetUserByUsername(ctx, username)

	if err != nil {
		// 如果错误是 ErrNoRows，说明没查到 -> 用户不存在
		if errors.Is(err, sql.ErrNoRows) {
			return nil // 用户不存在，返回 nil 表示没有错误
		}
		// 其他错误（数据库连接失败等） -> 向上抛出
		return err
	}
	return errors.New("用户已存在") // 用户存在，返回一个业务逻辑错误
}

func (l *UserLogic) SignUp(ctx context.Context, p *models.SignupParams) error {

	// 1. 判断用户是否存在
	if err := l.CheckUserExists(ctx, p.Username); err != nil {
		return err // 数据库查询出错或用户已存在
	}
	// 2. 生成UID
	userID := snowflake.GenID()

	// 假设你有这个工具函数
	hashedPassword, err := HashPassword(p.Password)
	if err != nil {
		return err
	}

	arg := db.CreateUserParams{
		UserID:   userID,
		Username: p.Username,
		Password: hashedPassword,
		Email:    p.Email,
		Gender:   false, // 默认性别，前端目前未提供该字段
	}

	// 4. 保存
	return l.store.CreateUser(ctx, arg)
}

func (l *UserLogic) Login(ctx context.Context, p *models.LoginParams) (user *models.User, err error) {
	// 1. 根据用户名查询用户
	userTemp, err := l.store.GetUserByUsername(ctx, p.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 2. 验证密码
	// 数据库中存储的是加密后的密码 user.Password
	// 登录参数中是明文密码 p.Password
	if err := CheckPassword(p.Password, userTemp.Password); err != nil {
		return nil, errors.New("密码错误")
	}

	user = &models.User{
		Username: userTemp.Username,
		UserID:   userTemp.UserID,
		Password: userTemp.Password,
	}

	// 3. 登录成功，生成 JWT token
	// 生成 token 的时候，我们需要用户的 ID 和用户名，这样我们在后续验证 token 时就能知道是哪个用户登录了

	//to do: 引入Redis，实现单端登录，生成新的 token 时将旧的 token 加入黑名单，过期时间设置为 token 的过期时间，在 JWTAuthMiddleware 中验证 token 时先检查是否在黑名单中，如果在黑名单中则拒绝访问
	//to do: 实现 token 刷新接口，客户端在 access token 过期后可以调用jwt.RefreshToken()来刷新access token，刷新时需要验证 refresh token 的有效性，如果 refresh token 有效且未过期，则生成新的 access token 和 refresh token，并返回给客户端，同时将旧的 refresh token 加入黑名单，过期时间设置为 refresh token 的过期时间

	user.Token, _, err = jwt.GenToken(user.UserID, user.Username)
	return

}
