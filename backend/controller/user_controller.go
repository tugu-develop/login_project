package controller

import (
	"fmt"
	"github.com/tugu-develop/login_project/model"
	"github.com/tugu-develop/login_project/usecase"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Authenticate(c *gin.Context) (string, error)
	AuthOk(c *gin.Context)
}

type userController struct {
	uu usecase.UserUsecase
}

func NewUserController(uu usecase.UserUsecase) UserController {
	return &userController{uu}
}

func (uc *userController) Signup(c *gin.Context) {
	user := model.User{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	err := uc.uu.Signup(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, nil)
}

func (uc *userController) Login(c *gin.Context) {
	user := model.User{}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	tokenString, err := uc.uu.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(1 * time.Hour),
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, tokenString)
}

func (uc *userController) Logout(c *gin.Context) {
	// リクエストのクッキーからトークンを取得
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, "Token not found")
		return
	}

	// ログアウト処理
	err = uc.uu.Logout(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, nil)
}

// 認証確認ミドルウェア
func (uc *userController) Authenticate(c *gin.Context) (string, error) {
	token, err := c.Cookie("token")
	if err != nil || token == "" { // トークン取得失敗
		fmt.Println("error: No token")
		return "", err
	}

	userID, err := uc.uu.Authenticate(token)
	if err != nil || userID == "" { // 認証失敗
		fmt.Println("error: No Authenticate")
		return "", err
	}

	return userID, nil
}

// ログイン認証OK
func (uc *userController) AuthOk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}
