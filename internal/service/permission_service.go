package service

import (
	"context"
	"math/rand"
	"testGolang/internal/dto"
	"time"
)

// PermissionsService севрис для работы с правами пользователя
type PermissionsService struct {
}

func NewPermissionsService() *PermissionsService {
	return &PermissionsService{}
}

// CheckAccess обращение к внешнему сервису, чтобы получить права пользователя
func (s *PermissionsService) CheckAccess(ctx context.Context, userID string) (*dto.PermissionsResponse, error) {
	duration := time.Duration(rand.Intn(70)) * time.Millisecond
	sleep, err := s.mockGetPermissionsUser(userID)
	select {
	case <-time.After(duration):
		return sleep, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// mockGetPermissionsUser мок прав пользователя
func (u *PermissionsService) mockGetPermissionsUser(userID string) (*dto.PermissionsResponse, error) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	check := r.Intn(2) == 1 // true или false
	return &dto.PermissionsResponse{
		CheckAccess: check,
	}, nil
}
