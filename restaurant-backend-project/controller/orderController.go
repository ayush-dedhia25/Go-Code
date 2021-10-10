package controller

import (
   "time"
   "fmt"
   "context"
   "net/http"
   
   "github.com/gin-gonic/gin"
   "gopkg.in/mgo.v2/bson"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/bson/primitive"
   
   "restaurant-backend-project/models"
   "restaurant-backend-project/database"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      result, err := orderCollection.Find(context.TODO(), bson.M{})
      if err != nil {
         msg := fmt.Sprintf("Error occurred while listing order items!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      var allOrders []models.Order
      if err = result.All(ctx, &allOrders); err != nil {
         log.Fatal(err)
      }
      c.JSON(http.StatusOK, allOrders)
   }
}

func GetOrder() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      orderId := c.Params("order_id")
      
      var order models.Order
      err := foodCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
      if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the order."})
      }
      c.JSON(http.StatusOK, order)
   }
}

func CreateOrder() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var table models.Table
      var order models.Order
      
      if err := c.BindJSON(&order); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      validationErr := validate.Struct(order)
      if validationErr != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
         return
      }
      
      if order.Table_id != nil {
         err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
         if err != nil {
            msg := fmt.Sprintf("Table was not found!")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
         }
      }
      
      order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      order.ID = primitive.NewObjectID()
      order.Order_id = food.ID.Hex()
      
      result, insertErr := orderCollection.InsertOne(ctx, order)
      if insertErr != nil {
         msg := fmt.Sprintf("Order item was not created!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
   }
}

func UpdateOrder() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var table models.Table
      var order models.Order
      
      orderId := c.Params("order_id")
      if err := c.BindJSON(&order); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      var updatedObj primitive.D
      
      if order.Table_id != nil {
         err := menuCollection.FindOne(ctx, bson.E{"table_id", food.Table_id}).Decode(&table)
         defer cancel()
         if err != nil {
            msg := fmt.Sprintf("Menu was not found!")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
         }
         updateObj = append(updateObj, bson.E{"menu": order.Table_id})
      }
      
      order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})
      
      upsert := true
      filter := bson.M{"order_id": orderId}
      opt := options.UpdateOptions{
         Upsert: &upsert,
      }
      
      result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj},}, &opt)
      if err != nil {
         msg := fmt.Sprintf("Order item update failed!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
   }
}

func OrderItemOrderCreator(order models.Order) string {
   order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
   order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
   order.ID = primitive.NewObjectID()
   order.Order_id = food.ID.Hex()
   
   orderCollection.InsertOne(ctx, order)
   defer cancel()
   
   return order.Order_id
}