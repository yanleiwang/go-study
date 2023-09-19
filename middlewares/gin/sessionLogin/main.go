package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
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
			//ctx.Redirect(http.StatusUnauthorized, "/login")
			// 不再继续执行业务 逻辑
			ctx.Abort()
			return
		}
		ctx.Next()

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
