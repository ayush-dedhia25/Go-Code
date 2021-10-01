package routes

import (
   "github.com/gin-gonic/gin"
   "github.com/ayush/restaurant-backend-project/controller"
)

func MenuRouter(incomingRoutes *gin.Engine) {
   incomingRoutes.GET("/menus", controller.GetMenus())
   incomingRoutes.GET("/menus/:menu_id", controller.GetMenu())
   incomingRoutes.POST("/menus", controller.CreateMenu())
   incomingRoutes.PATCH("/menus/:menu_id", controller.UpdateMenu())
}