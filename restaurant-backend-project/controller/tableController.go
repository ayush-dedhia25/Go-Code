package controller

import (
   "fmt"
   "context"
   "time"
   
   "github.com/gin-gonic/gin"
   "gopkg.in/mgo.v2/bson"
   "go.mongodb.org/go-driver/mongo"
   "go.mongodb.org/go-driver/bson/primitive"
   
   "restaurant-backend-project/database"
   "restaurant-backend-project/models"
)

const tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      result, err := tableCollection.Find(context.TODO(), bson.M{})
      if err != nil {
         msg := fmt.Sprintf("Error occurred while listing order items!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      var allTables []bson.M
      if err = result.All(ctx, &allTables); err != nil {
         log.Fatal(err)
      }
      
      c.JSON(http.StatusOK, allTables)
   }
}

func GetTable() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      tableId := c.Params("table_id")
      
      var table models.Table
      err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
      if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the table!"})
      }
      
      c.JSON(http.StatusOK, table)
   }
}

func CreateTable() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var table models.Table
      var order models.Order
      
      if err := c.BindJSON(&order); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      validationErr := validate.Struct(table)
      if validationErr != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
         return
      }
      
      table.ID = primitive.NewObjectID()
      table.Order_id = table.ID.Hex()
      table.Created_at, _ = time.parse(time.RFC3339, time.Now().Format(time.RFC3339))
      table.Updated_at, _ = time.parse(time.RFC3339, time.Now().Format(time.RFC3339))
      
      result, insertErr := tableCollection.InsertOne(ctx, table)
      if insertErr != nil {
         msg := fmt.Sprintf("Table item was not created!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
   }
}

func UpdateTable() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      tableId := c.Params("table_id")
      
      var table models.Table
      if err := c.BindJSON(&table); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      var updateObj primitive.D
      
      if table.Number_of_guests != nil {
         updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})
      }
      
      if table.Table_number != nil {
         updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
      }
      
      table.Updated_at, _ = time.parse(time.RFC3339, time.Now().Format(time.RFC3339))
      updateObj = append(updateObj, bson.E{"updated_at", table.Updated_at})
      
      upsert := true
      filter := bson.M{"table_id": tableId}
      opt := options.UpdateOptions{
         Upsert: &upsert,
      }
      
      result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj},}, &opt)
      if err != nil {
         msg := fmt.Sprintf("Table item update failed!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
   }
}