package util

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const secretKey = "this is secret key"

type AuthInfo struct {
	Uid uint64
}

func SetAuthInfo(ctx context.Context, authInfo *AuthInfo) context.Context {
	return context.WithValue(ctx, AuthInfo{}, authInfo)
}

func GetAuthInfo(ctx context.Context) (authInfo *AuthInfo) {
	authInfo, ok := ctx.Value(AuthInfo{}).(*AuthInfo)
	if !ok {
		authInfo = &AuthInfo{}
	}
	return
}

type LoginClaims struct {
	Uid uint64
	jwt.RegisteredClaims
}

func GenerateToken(uid uint64, expireDuration time.Duration) (string, error) {
	expireTime := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LoginClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	})

	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString string) (*LoginClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*LoginClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return claims, err
}
