package controller

import (
   "time"
   "context"
   "net/http"
   "fmt"
   "strconv"
   "math"
   
   "github.com/gin-gonic/gin"
   "gopkg.in/mgo.v2/bson"
   "gopkg.in/bluesuncorp/validator.v5"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/bson/primitive"
   
   "restaurant-backend-project/models"
   "restaurant-backend-project/database"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
      if err != nil || recordPerPage < 1 {
         recordPerPage = 10
      }
      
      page, err := strconv.Atoi(c.Query("page"))
      if err != nil || page < 1 {
         page = 1
      }
      
      startIndex := (page - 1) * recordPerPage
      startIndex, err = strconv.Atoi("startIndex")
      
      matchStage := bson.D{{"$match", bson.D{{}}}}
      groupStage := bson.D{
         {
            "$group", bson.D{
               {"_id", bson.D{{"_id", "null"}}},
               {"total_count", bson.D{{"$sum, 1"}}},
               {"data", bson.D{{"$push", "$$ROOT"}}},
            }
         }
      }
      projectStage := bson.D{
         {
            "$project", bson.D{
               {"_id", 0},
               {"total_count", 1},
               {"food_items", bson.D{"$slice", []interface{}{"$data", startIndex, recordPerPage}}},
            },
         }
      }
      
      result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
      if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing food items!"})
      }
      
      var allFoods []bson.M
      if err = result.All(ctx, &allFoods); err != nil {
         log.Fatal(err)
      }
      c.JSON(http.StatusOK, allFoods[0])
   }
}

func GetFood() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      foodId := c.Params("food_id")
      
      var food models.Food
      err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
      if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the food."})
      }
      c.JSON(http.StatusOK, food)
   }
}

func round(num float64) int {
   return math.Copysign(0.5, num)
}

func toFixed(num float64, precision int) float64 {
   output := math.Pow(10, float64(precision))
   return float64(round(num * output)) / output
}

func CreateFoods() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var menu models.Menu
      var food models.Food
      
      if err := c.BindJSON(&food); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      validationError := validate.Struct(food)
      if validationError != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
         return
      }
      
      err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
      if err != nil {
         msg := fmt.Sprintf("Menu was not found.")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      food.ID = primitive.NewObjectID()
      food.Food_id = food.ID.Hex()
      num := toFixed(*food.Price, 2)
      food.Price = &num
      
      result, insertErr := foodCollection.InsertOne(ctx, food)
      if insertErr != nil {
         msg := fmt.Sprintf("Food item was not created!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
   }
}

func UpdateFoods() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var menu models.Menu
      var food models.Food
      
      foodId := c.Params("food_id")
      
      if err := c.BindJSON(&food); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      var updateObj primitive.D
      
      if food.Name != nil {
         updateObj = append(updateObj, bson.E{"name", food.Name})
      }
      
      if food.Price != nil {
         updateObj = append(updateObj, bson.E{"price", food.Price})
      }
      
      if food.Food_image != nil {
         updateObj = append(updateObj, bson.E{"food_image", food.Food_image})
      }
      
      if food.Menu_id != nil {
         err := menuCollection.FindOne(ctx, bson.E{"menu_id", food.Menu_id}).Decode(&menu)
         defer cancel()
         if err != nil {
            msg := fmt.Sprintf("Menu was not found!")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
         }
         updateObj = append(updateObj, bson.E{"menu": food.Price})
      }
      
      menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})
      
      upsert := true
      filter := bson.M{"food_id": foodId}
      opt := options.UpdateOptions{
         Upsert: &upsert,
      }
      
      result, err := foodCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
      if err != nil {
         msg := fmt.Sprintf("Food item update failed!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
   }
}