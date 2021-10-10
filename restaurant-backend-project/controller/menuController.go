package controller

import (
   "time"
   "context"
   "log"
   "net/http"
   
   "github.com/gin-gonic/gin"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "gopkg.in/mgo.v2/bson"
   
   "restaurant-backend-project/database"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      result, err := menuCollection.Find(context.TODO(), bson.M{})
      if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing the menu items"})
      }
      
      var allMenus []bson.M
      if err = result.All(ctx, &allMenus); err != nil {
         log.Fatal(err)
      }
      c.JSON(http.StatusOK, allMenus)
   }
}

func GetMenu() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      menuId := c.Params("menu_id")
      
      var menu models.Menu
      err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
      if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the menu."})
      }
      
      c.JSON(http.StatusOK, menu)
   }
}

func CreateMenu() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var menu models.Menu
      if err := c.BindJSON(&menu); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      validationError := validate.Struct(menu)
      if validationError != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
         return
      }
      
      menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
      menu.ID = primitive.NewObjectID()
      menu.Menu_id = menu.ID.Hex()
      
      result, insertErr := menuCollection.InsertOne(ctx, menu)
      if insertErr != nil {
         msg := fmt.Sprintf("Menu was not created!")
         c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
         return
      }
      
      c.JSON(http.StatusOK, result)
      defer cancel()
   }
}

func inTimeSpan(start, end, check time.Time) bool {
   return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {
   return func(c *gin.Context) {
      ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
      defer cancel()
      
      var menu models.Menu
      if err := c.BindJSON(&menu); err != nil {
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      
      menuId := c.Params("menu_id")
      filter := bson.M{"menu_id": menuId}
      
      var updateObj primitive.D
      if menu.Start_Date != nil && menu.End_Date != nil {
         if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
            msg := "kindly retype the time"
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            defer cancel()
            return
         }
         
         updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
         updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})
         
         if menu.Name != nil {
            updateObj = append(updateObj, bson.E{"name", menu.Name})
         }
         
         if menu.Category != nil {
            updateObj = append(updateObj, bson.E{"category", menu.Category})
         }
         
         menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
         updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})
         
         upsert := true
         opt := options.UpdateOptions{
            Upsert: &upsert,
         }
         
         result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj},}, &opt)
         if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu update failed!"})
         }
         
         defer cancel()
         c.JSON(http.StatusOK, result)
      }
   }
}