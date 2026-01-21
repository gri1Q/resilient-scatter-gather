package service

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type VectorMemoryService struct {
}

func NewVectorMemoryService() *VectorMemoryService {
	return &VectorMemoryService{}
}

// GetContext время ответа от 100мс до 3 секунд
func (v *VectorMemoryService) GetContext(ctx context.Context, userID string) (string, error) {
	// симулируем случайную задержку 100ms - 3c
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	delay := time.Duration(100+r.Intn(2900)) * time.Millisecond

	select {
	case <-time.After(delay):
		// иногда падаем
		if r.Intn(10) == 0 {
			return "", errors.New("vector memory internal error")
		}
		return "context for user " + userID, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
