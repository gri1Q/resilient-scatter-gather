package service

import (
	"context"
	"math/rand"
	"testGolang/internal/dto"
	"time"
)

// UserService сервив для работы с пользователем
type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

// GetUser мок данные пользователя
// вызывает внешний сервис
func (s *UserService) GetUser(ctx context.Context, userID string) (*dto.UserResponse, error) {
	//случайную задержку от 10 до 50 мс
	sleep := time.Duration(10+rand.Intn(41)) * time.Millisecond
	select {
	case <-time.After(sleep):
		user := s.mockGetUser(userID)
		return user, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// mockGetUser мок данные пользователя
// имитация обращения во внешний севрис
func (u *UserService) mockGetUser(userID string) *dto.UserResponse {
	return &dto.UserResponse{
		ID:           userID,
		Name:         "User name",
		LastActivity: time.Now(),
	}
}
