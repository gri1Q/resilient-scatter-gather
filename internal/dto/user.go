package dto

import "time"

type UserResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	LastActivity time.Time `json:"last_activity"`
}
