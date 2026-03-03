package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("夏天夏天悄悄过去")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT token
func GenToken(userID int64, username string) (aToken, rToken string, err error) {
	// 创建一个新的 MyClaims 对象，包含用户 ID、用户名和标准 JWT 声明
	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置 token 的过期时间，这里设置为当前时间加上 TokenExpireDuration
			// jwt.NewNumericDate 是 jwt 包提供的一个函数，用于将 time.Time 转换为 *jwt.NumericDate 类型
			//time.Now() 返回当前时间，Add 方法用于在当前时间基础上添加一个时间段，这里是 TokenExpireDuration
			//Add方法传入的参数是一个 time.Duration 类型的值，表示要添加的时间长度，这里是 TokenExpireDuration，定义为 2 小时；返回一个新的 time.Time 对象，表示当前时间加上 2 小时后的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			//Issuer 是 JWT 的发行者，通常是一个字符串，表示这个 token 是由哪个系统或服务生成的，这里我们设置为 "my-auth-project"
			Issuer: "my-auth-project",
		},
	}
	// 生成一个新的 JWT token，使用 HS256 签名算法，并将 MyClaims 对象作为 token 的载荷
	// jwt.NewWithClaims 函数用于创建一个新的 JWT token，第一个参数是签名算法，这里我们使用 HS256，第二个参数是 token 的载荷，这里我们传入 MyClaims 对象 c，包含了用户 ID、用户名和标准 JWT 声明
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(MySecret)
	// 使用 MySecret 作为密钥对 token 进行签名，生成最终的 token 字符串
	// SignedString 方法会将 MySecret 作为密钥对 token 进行签名，并返回生成的 token 字符串和可能发生的错误

	//refresh token不需要存任何自定义数据，过期时间设置为7天
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		Issuer:    "my-auth-project",
	}).SignedString(MySecret)
	return aToken, rToken, err
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 创建一个新的 MyClaims 对象，用于存储解析后的 token 载荷
	// mc 是一个 MyClaims 类型的指针，用于存储解析后的 token 载荷，ParseWithClaims 方法会将解析后的 token 载荷填充到 mc 指向的 MyClaims 对象中
	var mc = new(MyClaims)
	// 解析 token 字符串，使用 MyClaims 作为 token 的载荷结构体
	// ParseWithClaims 方法传入 token 字符串、一个 MyClaims 对象的指针，以及一个回调函数用于提供签名验证所需的密钥
	// 解析 token 时，jwt 包会自动验证 token 的签名和过期时间，如果 token 无效或过期，ParseWithClaims 方法会返回一个错误

	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {

		return MySecret, nil
	})
	// 回调函数用于提供签名验证所需的密钥，这里我们直接返回 MySecret 作为密钥
	// 为什么可以直接用Mysecret作为密钥？因为我们在生成 token 时使用了 MySecret 进行签名，所以在解析 token 时也需要使用同样的密钥来验证签名的有效性
	if err != nil {
		return nil, err
	}

	// 验证 token 是否有效，token.Valid 是 jwt 包提供的一个字段，用于表示 token 是否有效，如果 token 无效或过期，token.Valid 将为 false
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Access Token
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	//refresh token无效直接返回错误
	if _, err := ParseToken(rToken); err != nil {
		return "", "", err
	}

	var claims = new(MyClaims)
	//解析 Access Token，判断是否过期
	token, err := jwt.ParseWithClaims(aToken, claims, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})

	//如果 Access Token 是过期错误，且 Refresh Token 有效，则生成新的 Access Token 和 Refresh Token
	//ValidationError 是 jwt 包提供的一个错误类型，用于表示 token 验证过程中发生的错误，如果 err 是一个 ValidationError 类型的错误，并且错误类型包含 ValidationErrorExpired，说明 Access Token 已经过期了，这时我们可以使用 Refresh Token 来生成新的 Access Token 和 Refresh Token
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return GenToken(claims.UserID, claims.Username)
		}

	}
	//其他错误直接返回
	if err != nil {
		return "", "", err
	}
	//如果 Access Token 仍然有效，则直接返回原来的 Access Token 和 Refresh Token
	if token.Valid {
		return aToken, rToken, nil
	}
	return "", "", errors.New("invalid token")
}
