package main

import "github.com/gin-gonic/gin"

func main() {
	server := gin.Default()

	users := server.Group("/users")
	orders := server.Group("/orders")

	// /orders/:id
	orders.GET("/:id", func(c *gin.Context) {
		oid := c.Param("id")
		c.JSON(200, gin.H{
			"message": "这是订单信息，oid 是" + oid,
		})
	})

	vipGroup := users.Group("/vip")
	// /users/vip/upgrade
	vipGroup.GET("/upgrade", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "大冤种，你要升级成 VVIP 吗",
		})
	})

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
