package models

import (
   "time"
   jwt "github.com/dgrijalva/jwt-go"
   "secure-api-project/randomstrings"
)

const (
   RefreshTokenValidTime = time.Hour * 72
   AuthTokenValidTime = time.Minute * 15
)

type User struct {
   Username     string
   PasswordHash string
   Role         string
}

type TokenClaims struct {
   jwt.StandardClaims
   Role string `json:"role"`
   Csrf string `json:"csrf"`
}

func GenerateCsrfSecret() (string, error) {
   return randomstrings.GenerateRandomString(32)
}