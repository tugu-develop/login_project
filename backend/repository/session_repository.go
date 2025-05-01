package repository

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, jti string, userID string) error
	GetUserIDByToken(ctx context.Context, token string) (string, error)
	DeleteSession(ctx context.Context, token string) error
}

type redisSessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) SessionRepository {
	return &redisSessionRepository{client: client}
}

// ログイン時にセッションを保存
func (r *redisSessionRepository) SaveSession(ctx context.Context, jti string, userID string) error {
	return r.client.Set(ctx, jti, userID, time.Hour*12).Err()
}

// API リクエスト時に Redis からユーザーIDを取得
func (r *redisSessionRepository) GetUserIDByToken(ctx context.Context, jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// 署名の検証
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// claimsを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	// jtiを取り出す
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return "", fmt.Errorf("jti not found in token")
	}

	// Redisから jti が一致するセッション情報を取得
	redisKey := "session:jti:" + jti
	userID, err := r.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("session not found for jti")
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

// ログアウト時にセッションを削除
func (r *redisSessionRepository) DeleteSession(ctx context.Context, jwtToken string) error {
	// JWTトークンの解析
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// 署名の検証
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token: %w", err)
	}

	// claimsを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid claims")
	}

	// jtiを取り出す
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return fmt.Errorf("jti not found in token")
	}

	// Redisから jti に基づくセッション情報を削除
	redisKey := "session:jti:" + jti
	err = r.client.Del(ctx, redisKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}