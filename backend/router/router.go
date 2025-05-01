package router

import (
	"login_project/controller"
	"login_project/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(uc controller.UserController) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"DELETE",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// ヘルスチェック用
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/signup", uc.Signup)
	r.POST("/login", uc.Login)
	r.POST("/logout", uc.Logout)

	// 認証が必要なエンドポイント（/auth）
	authGroup := r.Group("/auth")  // ここを変更
	authGroup.Use(middleware.AuthMiddleware(uc))
	{
		authGroup.GET("", uc.AuthOk)  // ハンドラーを設定
	}

	return r
}
