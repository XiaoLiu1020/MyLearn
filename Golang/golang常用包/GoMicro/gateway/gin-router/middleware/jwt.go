package middleware

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("TodoList")

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code uint32

		code = 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			_, err := ParseToken(token)
			if err != nil {
				code = 401
			}
		}
		if code != 200 {
			c.JSON(500, gin.H{
				"code": code,
				"msg":  "鉴权失败",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

type Claims struct {
	Id uint `json:"id"`
	jwt.StandardClaims
}

// 签发用户token
func GenerateToken(id uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "todoList",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 验证用户token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
