package main

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const secretKey = "this is secret key"

type LoginClaims struct {
	PhoneNumber string
	jwt.RegisteredClaims
}

func GenerateToken(phoneNumber string, expireDuration time.Duration) (string, error) {
	expireTime := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LoginClaims{
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "wolf",
		},
	})
	// SecretKey 用于对用户数据进行签名，不能暴露
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Printf("verify token failed: %s\n", err)
		return false
	}

	claims, ok := token.Claims.(*LoginClaims)
	if ok && token.Valid {
		fmt.Printf("%v %v %v\n", claims.PhoneNumber, claims.RegisteredClaims.Issuer, claims.RegisteredClaims.ExpiresAt)
		return true
	}

	fmt.Printf("verify token failed: %s\n", err)
	return false
}

func main() {
	tokenString, _ := GenerateToken("18820137896", 3*time.Second)

	fmt.Printf("generated token string: %v\n", tokenString)

	if VerifyToken(tokenString) {
		fmt.Println("verify token succ")
	} else {
		fmt.Println("verify token failed")
	}

	time.Sleep(4 * time.Second)

	if VerifyToken(tokenString) {
		fmt.Println("verify token succ")
	} else {
		fmt.Println("verify token failed")
	}
}
