package jwttoken

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
)

// Sign 加密Token
func Sign(s interface{}, expires time.Duration, tokenKey, tokenSecret string) (string, error) {
	now := time.Now()
	of := reflect.TypeOf(s)
	var ss string
	switch of.String() {
	case "string":
		ss = s.(string)
	case "int":
		ss = strconv.Itoa(s.(int))
	}
	expiresAt := now.Add(expires * time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Id:        ss,
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    tokenKey,
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", errors.Wrap(err, "sign err")
	}
	return tokenString, nil
}

// Unsign 解密Token
func Unsign(tokenString, tokenSecret string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "token invalid")
	}
	s := token.Claims.(*jwt.StandardClaims).Id
	if s == "" {
		return nil, errors.Wrap(err, "unsign err")
	}
	return s, nil
}
