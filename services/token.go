package services

import (
	"KnowEase/models"
	"KnowEase/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

// 生成JWT token
func (ts *TokenService) GenerateToken(user *models.User) string {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 180).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte("secret"))

	return t
}

// 检验token
func (ts *TokenService) VerifyToken(tokenString string) (string, error) {
	if ts.IsTokenBlacklisted(tokenString) {
		return "", fmt.Errorf("token has been invalidated")
	}
	if tokenString == "" {
		return "", errors.New("authorization token required")
	}

	// 去掉 Bearer 前缀
	if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}

	// 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	//检验是否有效
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role not found in token")
	}
	return role, nil
}

// 将token加入黑名单
func (ts *TokenService) AddToBlacklist(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	expiration, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("token does not have exp claim")
	}

	err = utils.Client.Set(utils.Ctx, tokenString, "blacklisted", time.Unix(int64(expiration), 0).Sub(time.Now())).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %v", err)
	}

	return nil
}

// 检测token是否在黑名单内
func (ts *TokenService)IsTokenBlacklisted(tokenString string) bool {
	val, err := utils.Client.Get(utils.Ctx, tokenString).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		return false
	}
	return val == "blacklisted"
}
