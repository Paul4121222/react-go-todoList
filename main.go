package main
import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"log"
)

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Completed bool `bson:"completed" json:"completed"`
	Message string `bson:"body" json:"body"`
}

var collection *mongo.Collection;

func handleGetTodos (c *gin.Context) {
	var todos []Todo;
	cursor, err := collection.Find(context.Background(), bson.M{});
	if(err != nil) {
		c.JSON(500, gin.H{
				"error": err,
		});
		log.Fatal("搜尋失敗")
	}

	//應釋放資源
	defer cursor.Close(context.Background());

	for cursor.Next(context.Background()) {
		var todo Todo;
		if err := cursor.Decode(&todo); err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		} else {
			todos = append(todos, todo);
		}
	}
	c.JSON(200, todos)
}

func handleAddTodo(c *gin.Context) {
	todo := Todo{};
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(500, gin.H{
			"error": err,
		});
		return;
	} 
	if todo.Message == "" {
		c.JSON(400, gin.H{
			"error": "Body is Required",
		})
		return;
	}	

	result, err := collection.InsertOne(context.Background(), todo);
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		});
		return;
	}

 	c.JSON(200, result);
}	

func handleUpdateTodo(c *gin.Context) {
	 id := c.Param("id");
	 objectID, err := primitive.ObjectIDFromHex(id);

	 if (err != nil) {
		c.JSON(400, gin.H{
			"error": "ID is Invalid.",
		})
		return;
	 }
	 result, err := collection.UpdateOne(
		context.Background(), 
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"completed": true}},
	);

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		});
		return;
	}

 	c.JSON(200, result);
}

func handleDeleteTodo(c *gin.Context) {
	id := c.Param("id");
	objectID, err := primitive.ObjectIDFromHex(id);
	if (err != nil) {
		c.JSON(400, gin.H{
			"error": "ID is Invalid.",
		})
		return;
	 }

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID});
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		});
		return;
	}

 	c.JSON(200, result);
}
func main() {
	env := os.Getenv("ENV");

	//生產環境下不須加載env
	if env != "production" {
		if err := godotenv.Load(); err != nil {
				fmt.Println("Not load .env file.");
				return;
		}
	}
	

	port := os.Getenv("PORT");
	mongodbURI := os.Getenv("MONGO_URI");
	clientOptions := options.Client().ApplyURI(mongodbURI);

	//跟mongoDB做連線(去問ip是否有效，比方說給的url無法做dns轉換就會發生err)，並返回客戶端的連線管理器
	client, err := mongo.Connect(context.Background(), clientOptions);
	if err != nil {
		log.Fatal("無法初始化 Mongo 客戶端:", err);
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("無法連線到 MongoDB:", err)
	}

	defer client.Disconnect(context.Background());
	fmt.Println("成功連線到mongoDB");
	collection = client.Database("go-react").Collection("todos");

	r := gin.Default()
	
	r.GET("/api/todos", handleGetTodos);
	r.POST("/api/todo", handleAddTodo);
	r.PATCH("/api/todo/:id", handleUpdateTodo);
	r.DELETE("/api/todo/:id", handleDeleteTodo);

	
	if env == "production" {
		//vite 會將靜態資源打包在/assets目錄下
		r.Static("/assets", "./client/dist/assets");
		r.NoRoute(func (c *gin.Context) {
		c.File("./client/dist/index.html")
	})
	} else {
		//允許跨域
		r.Use(cors.Default());
	}


	r.Run(":" + port);
}