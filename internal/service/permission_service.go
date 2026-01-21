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

	sleep, err := s.mockGetPermissionsUser(userID)
	select {
	case <-time.After(50 * time.Millisecond):
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
