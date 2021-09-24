package db

import (
   "errors"
   "log"
   "golang.org/x/crypto/bcrypt"
   "secure-api-project/db/models"
)

var users = map[string]models.User{}

func InitDB() {}

func FetchUserByUsername(username string) (models.User, string, error) {
   for k, v := range users {
      if v.Username == username {
         return v, k, nil
      }
   }
   return models.User{}, "", errors.New("User not found that matches the username!")
}

func StoreUser(username, password, role string) (uuid string, err error) {}

func DeleteUser() {}

func FetchUserById() {}

func StoreRefreshToken() {}

func DeleteRefreshToken() {}

func CheckRefreshToken() bool {}

func LogUserIn() {}

func generateBcryptHash() {}

func checkPasswordAgainstHash() error {}