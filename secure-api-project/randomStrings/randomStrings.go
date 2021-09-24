package randomstrings

import (
   "crypto/rand"
   "encoding/base64"
)

func generateRandomBytes(size int) ([]byte, error) {
   b := make([]byte, size)
   _, err := rand.Read(b)
   if err != nil {
      return nil, err
   }
   return b, nil
}

func GenerateRandomString(size int) (string, error) {
   byteString, err := generateRandomBytes(size)
   return base64.URLEncoding.EncodeToString(byteString), err
}