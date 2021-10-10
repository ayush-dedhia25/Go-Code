package routes

import (
   "github.com/gin-gonic/gin"
   "restaurant-backend-project/controller"
)

func TableRouter(incomingRoutes *gin.Engine) {
   incomingRoutes.GET("/tables", controller.GetTables())
   incomingRoutes.GET("/tables/:table_id", controller.GetTable())
   incomingRoutes.POST("/tables", controller.CreateTable())
   incomingRoutes.PATCH("/tables/:table_id", controller.UpdateTable())
}