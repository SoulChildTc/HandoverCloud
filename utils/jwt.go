package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"soul/global"
	"time"
)

var (
	ErrTokenSignatureInvalid = errors.New("token signature method is invalid")
	ErrTokenIllegal          = errors.New("token is illegal")
)

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func CreateJwtToken(userid int, username string) (string, error) {
	signingKey := []byte(global.Config.Jwt.Secret)
	ttl := global.Config.Jwt.Ttl
	iss := global.Config.AppName

	claims := CustomClaims{
		UserID:   userid,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(ttl)),
			NotBefore: jwt.NewNumericDate(time.Now().Local().Add(-time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return ss, nil

}

func ParseJwtToken(tokenString string) (*CustomClaims, error) {
	fmt.Println(tokenString)
	signingKey := []byte(global.Config.Jwt.Secret)
	tokenObj, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenObj.Method != jwt.SigningMethodHS256 {
		return nil, ErrTokenSignatureInvalid
	}

	claims, ok := tokenObj.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrTokenIllegal
	}

	return claims, nil
}
