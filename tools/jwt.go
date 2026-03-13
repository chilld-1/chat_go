package tools

import (
	"errors"
	"gochat/db"
	"gochat/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("sxl")

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToke(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
func VerifyToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	return claims, err
}

func TokenCheck(c *gin.Context, tokenString string) bool {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		if c != nil {
			UnauthorizedResponse(c, "无效的认证令牌: "+err.Error())
		}
		return false
	}
	err = db.DB.Where("username=?", claims.Username).First(&model.User{}).Error
	if err != nil {
		if c != nil {
			UnauthorizedResponse(c, "用户不存在")
		}
		return false
	}

	return true
}
