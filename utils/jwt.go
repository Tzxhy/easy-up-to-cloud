package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type UserTokenClaims struct {
	Username string `json:"username"`
	UserId   string `json:"userid"`
	jwt.StandardClaims
}

var MySecret = []byte("qva5im03q96fnjaga1rnafp3qrsi8r")

const TokenExpireDuration = time.Hour * 24

// GenUserToken 生成JWT
func GenUserToken(username string, userid string) (string, error) {
	// 创建一个我们自己的声明
	c := UserTokenClaims{
		username, // 自定义字段
		userid,
		jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(TokenExpireDuration).Unix())), // 过期时间
			Issuer:    "root",                                                           // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseUserToken 解析JWT
func ParseUserToken(tokenString string) (*UserTokenClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserTokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserTokenClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

type ShareItemClaims struct {
	Sid string `json:"sid"`
	jwt.StandardClaims
}

// GenUserToken 生成JWT
func GenShareItemToken(sid string) (string, error) {
	// 创建一个我们自己的声明
	c := ShareItemClaims{
		sid, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(TokenExpireDuration).Unix())), // 过期时间
			Issuer:    "root",                                                           // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseUserToken 解析JWT
func ParseShareItemToken(tokenString string) (*ShareItemClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &ShareItemClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*ShareItemClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
