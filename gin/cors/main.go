package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

/*
1. 运行从测试， 访问  http://127.0.0.1:8081/
2. 点击跨域请求按钮， 会访问 'http://127.0.0.1:9091/data'

查看http响应 可以看到， 默认情况下， cors设置了以下4个http 响应头：

Access-Control-Allow-Headers: Origin,Content-Length,Content-Type
Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS
Access-Control-Allow-Origin: *
Access-Control-Max-Age: 43200

Access-Control-Allow-Headers： 跨域请求支持的 http head
Access-Control-Allow-Methods: 跨站请求支持的 http请求方法
Access-Control-Allow-Origin：  支持哪些域名跨站请求， * 表示任意域名
Access-Control-Max-Age：用来指定本次预检请求的有效期，

所以如果自己实现cors， 也可以简单粗暴的 设置这些head就行
*/

func main() {
	go func() {
		s := gin.Default()
		s.LoadHTMLGlob("*.html")
		s.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		if err := s.Run(":8081"); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second * 2)
	s := gin.Default()

	// 允许跨站请求， 注释掉  请求就不会被发出来。
	s.Use(cors.New(cors.Config{
		AllowCredentials: true,                                                 // 是否允许带 cookie
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"}, // 允许的请求头
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour, // preflight响应 过期时间
	}))

	s.PUT("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "进来了这里",
		})
	})
	if err := s.Run(":9091"); err != nil {
		panic(err)
	}
}
