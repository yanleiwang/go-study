package main

import "github.com/gin-gonic/gin"

func main() {
	server := gin.Default()
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.POST("/post", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "这是一个 POST 方法",
		})
		c.Bind()
	})

	server.GET("/users/:name", func(c *gin.Context) {
		username := c.Param("name")
		c.JSON(200, gin.H{
			"message": "这是用户信息，uid 是" + username,
		})
	})

	// /order?id=123
	server.GET("/order", func(c *gin.Context) {
		oid := c.Query("id")
		c.JSON(200, gin.H{
			"message": "这是订单信息，oid 是" + oid,
		})
	})

	// 通配符路由 只能  /xx/*xxx, 其他都不行， 比如： 不能注册这种 /users/*, /users/*/a
	server.GET("/views/*.html", func(c *gin.Context) {
		name := c.Param(".html")
		c.JSON(200, gin.H{
			"message": "这是页面，文件名是" + name,
		})
	})

	//server.GET("/viewsv1/*", func(c *gin.Context) {
	//	name := c.Param(".html")
	//	c.JSON(200, gin.H{
	//		"message": "这是页面，文件名是" + name,
	//	})
	//})

	// ip:port
	// 监听并在 0.0.0.0:8080 上启动服务
	if err := server.Run(); err != nil {
		panic(err)
	}
}
