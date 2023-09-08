package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("./templates/*.gohtml")

	server.GET("/form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.gohtml", nil)
	})

	server.POST("/form", func(c *gin.Context) {
		var u User
		c.Bind(&u)
		fmt.Printf("user_name: %s \n password: %s \n", u.Name, u.Password)
		c.JSON(200, gin.H{
			"message": "收到了数据" + u.Name,
		})
	})

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}

type User struct {
	Name     string `form:"user_name"`
	Password string `form:"password"`
}
