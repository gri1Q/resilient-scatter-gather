package service

import (
	"context"
	"math/rand"
	"testGolang/internal/apperrors"
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
func (s *UserService) GetUser(ctx context.Context, userID string) (*dto.UserResponse, error) {
	//случайную задержку от 10 до 14
	duration := time.Duration(rand.Intn(14)) * time.Millisecond

	timer := time.NewTimer(duration)
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	select {
	case <-timer.C:
		user := s.mockGetUser(userID)
		return user, nil
	case <-ctx.Done():
		return nil, apperrors.ErrTimeout
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
