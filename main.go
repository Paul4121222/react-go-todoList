package main
import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Todo struct {
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Message string `json:"body"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Not load .env file.");
		return;
	}

	port := os.Getenv("PORT");

	r := gin.Default();
	todoList := []Todo{};

	r.GET("/api/todos", func (c *gin.Context) {
		c.JSON(200, todoList)
	})


	r.POST("/api/todo", func (c *gin.Context){
		todo := Todo{}
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid",
			})
			return;
		}

		if todo.Message == "" {
			c.JSON(400, gin.H{
				"message": "Body is Required",
			})
			return;
		}
		todo.ID = len(todoList) + 1;

		//一定要返回接收值，原因是當slice容量不構，記憶體會分配更大的cap，再把舊/新資料搬過去，然後返回新slice
		todoList = append(todoList, todo);

		c.JSON(200, todo)
	})

	r.PATCH("/api/todo/:id", func(c *gin.Context) {
		id := c.Param("id");

		for i, todo := range todoList {
			if fmt.Sprint(todo.ID) == id {
				todoList[i].Completed = true;
				c.JSON(200, todo);
				return;
			}
		}

		c.JSON(404, gin.H{
			"message": "Not found.",
		})
	})

	r.DELETE("/api/todo/:id", func(c *gin.Context) {
		id := c.Param("id");
		for i, todo := range todoList {
			if(fmt.Sprint(todo.ID) == id) {
				todoList = append(todoList[:i], todoList[(i+1):]...)
				c.JSON(200, gin.H{
					"message": "Success",
				})
				return;
			}
		}

		c.JSON(404, gin.H{
			"message": "Not found.",
		})
	})

	r.Run(":" + port)
}