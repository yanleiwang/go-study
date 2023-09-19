package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// NewStore 最后参数 实际是 hashKey和blockKey（可省略），前者用于验证，后者用于加密
	// hashkey 用于 给 session hash值进行加密， 用于保证 session 没有被篡改
	// blockkey 用于 给 session 数据进行加密
	// []byte("secret") 就是hashKey
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	// mysession， 指定了cookie里的 key 名称
	r.Use(sessions.Sessions("mysession", store))

	// 某个用户每次访问该页面， 就给 count++
	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8000")
}
