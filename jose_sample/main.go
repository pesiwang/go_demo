package main

import (
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func main() {
	// 创建 JWT 签名器
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       []byte("secret"),
	}, nil)
	if err != nil {
		panic(err)
	}

	// 创建 JWT 的载荷
	claims := jwt.Claims{
		Subject:   "user-id-123",
		Issuer:    "http://example.com",
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	// 创建 JWT
	rawToken, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		panic(err)
	}
	fmt.Println(rawToken)

	// 验证 JWT
	verificationKey := jose.VerificationKey{
		Algorithm: jose.HS256,
		Key:       []byte("secret"),
	}
	verifier, err := jose.NewVerifier(verificationKey)
	if err != nil {
		panic(err)
	}
	parser, err := jwt.ParseSigned(rawToken)
	if err != nil {
		panic(err)
	}
	if err := parser.Verify(verifier); err != nil {
		panic(err)
	}
	var claims jwt.Claims
	if err := parser.Claims(&claims); err != nil {
		panic(err)
	}
	fmt.Println(claims)
}
