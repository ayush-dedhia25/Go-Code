package controller

import (
   "fmt"
   "log"
   "github.com/gin-gonic/gin"
   "go.mongodb.org/mongo-driver/mongo"
)

func GetUsers() gin.HandlerFunc {
   return func(c *gin.Context) {
      
   }
}

func GetUser() gin.HandlerFunc {
   return func(c *gin.Context) {
      
   }
}

func Signup() gin.HandlerFunc {
   return func(c *gin.Context) {
      
   }
}

func Login() gin.HandlerFunc {
   return func(c *gin.Context) {
      
   }
}

func HashPassword(password string) string {
   
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
   
}