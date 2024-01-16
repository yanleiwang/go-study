package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	// NewStore 最后参数 实际是 hashKey和blockKey（可省略），前者用于验证，后者用于加密
	// hashkey 用于 给 session hash值进行加密， 用于保证 session里 的内容 没有被篡改
	// blockkey 用于 给 session 数据进行加密
	// 推荐使用32位
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "",
		[]byte("pY8tX3vY7aT8nK2nD6lO9jR4pE5aN4gI"), []byte("rM8eL5rB7pC1fZ4tZ3eT1fM8cS5kK7lD"))
	// mysession， 指定了cookie里的 key 名称
	r.Use(sessions.Sessions("mysession", store))

	// 请求过来的时候， 校验有没有登录
	r.Use(func(ctx *gin.Context) {
		// 是登录， 直接执行业务逻辑（验证用户名， 密码）
		// 否则跳转到登录页面
		if ctx.FullPath() == "/login" {
			return
		}

		// 不是登录， 就要求 cookie 带上合法的session_id， 能够找到对应的session
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.String(http.StatusUnauthorized, "你没登录哦")
			ctx.Abort()
			return
		}
		// 每隔一段时间 刷新过期时间
		const timeKey = "updateTime"
		val := sess.Get(timeKey)
		now := time.Now().UnixMilli()
		updateTime, ok := val.(int64)
		// 处于演示效果，整个 session 的过期时间是 1 分钟，所以我这里十秒钟刷新一次。
		// val == nil 是说明刚登录成功
		// 我们不在登录里面初始化这个 update_time，是因为它属于"刷新"机制，而不属于登录机制
		if val == nil || (ok && now-updateTime > 10*1000) {
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			sess.Set(timeKey, now)
			err := sess.Save()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		if val != nil && !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	})

	// localhost:8000/login?name=qqq&password=123
	r.GET("/login", func(ctx *gin.Context) {
		name := ctx.Query("name")
		password := ctx.Query("password")

		// 用户名， 密码验证通过
		if name == "qqq" && password == "123" {
			// 设置session
			sess := sessions.Default(ctx)
			sess.Set("userId", 123)

			// Gin 的 Session 用这些选项来初始化 Cookie。除了 MaxAge 有多层含义，其它参数就是在 Cookie里的含义
			//MaxAge 则不同，它一方面用来控制 Cookie，而有一些实现，也用它来控制 Session 中的 key、value 的过
			//期时间。
			//比如 Redis，它会用这个来控制你的数据的过期时间
			sess.Options(sessions.Options{
				Secure:   true,
				HttpOnly: true,
				// 一分钟过期
				MaxAge: 60,
			})
			err := sess.Save() //  不save 不会真的写入redis
			if err != nil {
				ctx.String(http.StatusInternalServerError, "系统错误")
				return
			}
			ctx.String(http.StatusOK, "登录成功")
			return
		}

		ctx.String(http.StatusOK, "用户名， 密码错误了")
	})

	r.GET("/profile", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你登录了， 这是profile")
	})

	r.Run(":8000")
}
