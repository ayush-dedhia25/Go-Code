package main

import (
   "os"
   "github.com/gin-gonic/gin"
   "go.mongodb.org/mongo-driver/mongo"
   "github.com/ayush/restaurant-backend-project/database"
   "github.com/ayush/restaurant-backend-project/routes"
   "github.com/ayush/restaurant-backend-project/middleware"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
   port := os.Getenv("PORT")
   if port == "" {
      port = "8080"
   }
   
   router := gin.New()
   router.Use(gin.Logger())
   router.Use(middleware.Authentication)
   
   routes.FoodRoutes(router)
   routes.InvoiceRoutes(router)
   routes.MenuRoutes(router)
   routes.OrderItemRoutes(router)
   routes.OrderRoutes(router)
   routes.TableRoutes(router)
   routes.UserRoutes(router)
   
   router.Run(":" + port)
}