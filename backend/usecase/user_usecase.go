package usecase

import (
	"context"
	"fmt"
	"login_project/model"
	"login_project/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Signup(user model.User) error
	Login(user model.User) (string, error)
	Logout(token string) error
	Authenticate(token string) (string, error)
}

type userUsecase struct {
	ur repository.UserRepository
	sr repository.SessionRepository
}

func NewUserUsecase(ur repository.UserRepository, sr repository.SessionRepository) UserUsecase {
	return &userUsecase{
		ur: ur,
		sr: sr,
	}
}

func (uu *userUsecase) Signup(user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return err
	}
	return nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// UUID を生成
	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"jti":     jti,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	fmt.Println("=========================")
	fmt.Println(token)
	fmt.Println("=========================")
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	// Redis にセッション保存（キー: トークン, 値: ユーザーID）
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session:jti:%s", jti)  // Redisキーsession:jti:<jti>
	if err := uu.sr.SaveSession(ctx, sessionKey, fmt.Sprintf("%v", storedUser.ID)); err != nil {
		return "", err
	}

	return tokenString, nil
}

// ログアウト時に Redis から削除
func (uu *userUsecase) Logout(token string) error {
	ctx := context.Background()
	return uu.sr.DeleteSession(ctx, token)
}

// 認証チェック（Redis でセッション確認）
func (uu *userUsecase) Authenticate(token string) (string, error) {
	ctx := context.Background()
	userID, err := uu.sr.GetUserIDByToken(ctx, token)
	if err != nil {
		return "", err
	}
	if userID == "" {
		return "", nil // 未認証
	}
	return userID, nil
}
