package main

import (
   "log"
   "secure-api-project/db"
   "secure-api-project/server"
   "secure-api-project/myJWT"
)

const (
   HOST = "localhost"
   PORT = "8000"
)

func main() {
   db.InitDB()
   jwtError := myJWT.InitJWT()
   
   if jwtError != nil {
      log.Println("Error initializing JWT")
      log.Fatal(jwtError)
   }
   
   serverError := server.StartServer(HOST, PORT)
   if serverError != nil {
      log.Println("Error starting server!")
      log.Fatal(serverError)
   }
}