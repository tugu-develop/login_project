package middleware

import (
	"login_project/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(uc controller.UserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証処理
		userID, err := uc.Authenticate(c)
		if err != nil || userID == "" { // 認証NG
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"authenticated": false})
			return
		}
		// ユーザーIDをコンテキストに保存（後続のハンドラで使えるかも）
		c.Set("user_id", userID)
		c.Next()
	}
}
