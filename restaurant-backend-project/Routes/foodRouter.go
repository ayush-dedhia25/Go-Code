package routes

import (
   "github.com/gin-gonic/gin"
   "github.com/ayush/restaurant-backend-project/controller"
)

func FoodRouter(incomingRoutes *gin.Engine) {
   incomingRoutes.GET("/foods", controller.GetFoods())
   incomingRoutes.GET("/foods/:food_id", controller.GetFood())
   incomingRoutes.POST("/foods", controller.CreateFood())
   incomingRoutes.PATCH("/foods/:food_id", controller.UpdateFood())
}