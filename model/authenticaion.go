package model

import (
	"github.com/dgrijalva/jwt-go"
)

// LoginInfo model
type LoginInfo struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password" `
}

// JWTCustomClaims model
type JWTCustomClaims struct {
	Email     string `json:"email"`
	TokenType string `json:"tokenType"`
	Admin     bool   `json:"admin"`
	jwt.StandardClaims
}
