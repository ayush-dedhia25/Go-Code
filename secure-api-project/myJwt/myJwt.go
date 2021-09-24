package myJWT

import (
   "crypto/rsa"
   "io/ioutil"
   "time"
   "errors"
   "log"
   jwt "github.com/dgrijalva/jwt-go"
   "secure-api-project/db"
   "secure-api-project/db/models"
)

const (
   privateKeyPath = "keys/app.rsa"
   publicKeyPath = "keys/app.rsa.pub"
)

func InitJWT() error {
   signBytes, err := ioutil.ReadFile(privateKeyPath)
   if err != nil {
      return err
   }
   
   signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
   if err != nil {
      return err
   }
   
   verifyBytes, err := ioutil.ReadFile(publicKeyPath)
   if err != nil {
      return err
   }
   
   verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
   if err != nil {
      return err
   }
   return nil
}

func createNewTokens(uuid, role string) (authTokenString, refreshTokenString, csrfSecret string, err error) {
   // Generating CSRF Secret
   csrfSecret, err := models.GenerateCsrfSecret()
   if err != nil {
      return
   }
   
   // Generating refresh token
   refreshTokenString, err = createRefreshTokenString(uuid, role, csrfSecret)
   if err != nil {
      return 
   }
   
   // Generating auth token
   authTokenString, err = createAuthTokenString(uuid, role, csrfSecret)
   if err != nil {
      return
   }
   return
}

func CheckAndRefreshTokens() () {}

func createAuthTokenString(uuid, role, csrfSecret string) (authTokenString string, err error) {
   authTokenExp := time.Now().Add(models.AuthTokenValidTime).Unix()
   authClaims := models.TokenClaims{
      jwt.StandardClaims{
         Subject: uuid,
         ExpiresAt: authTokenExp,
      },
      Role: role,
      Csrf: csrfSecret,
   }
   authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS265"), authClaims)
   authTokenString, err := authJwt.SignedString(signKey)
   if err != nil {
      return err
   }
   return
}

func createRefreshTokenString(uuid, role, csrfSecret string) (refreshTokenString string, err error) {
   refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()
   refreshJti, err := db.StoreRefreshToken()
   if err != nil {
      return
   }
   refreshClaims := models.TokenClaims{
      jwt.StandardClaims{
         Id: refreshJti,
         Subject: uuid,
         ExpiresAt: refreshTokenExp,
      },
      Role: role,
      Csrf: csrfSecret,
   }
   refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS265"), refreshClaims)
   refreshTokenString, err = refreshJwt.SignedString(signKey)
   if err != nil {
      return
   }
   return
}

func updateRefreshTokenExp() () {
   
}

func updateAuthTokenString() () {
   
}

func RevokeRefreshToken() error {
   
}

func updateRefreshTokenCsrf() () {
   
}

func grabUUID() () {
   
}