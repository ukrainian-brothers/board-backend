package user

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type User struct{}

type Social struct {
	UserID       uuid.UUID       `json:"user_id"`
	Social       string          `json:"social"`
	SocialId     string          `json:"social_id"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	UserData     json.RawMessage `json:"user_data"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DestroyedAt  time.Time       `json:"destroyed_at"`
}
