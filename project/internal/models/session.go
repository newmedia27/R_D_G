package models

import "time"

type Session struct {
	ID               string    `json:"id,omitempty"`
	UserID           string    `json:"user_id"`
	RefreshTokenHash string    `json:"-"`
	DeviceSessionID  string    `json:"device_session_id"`
	FingerPrint      string    `json:"-"`
	UserAgent        string    `json:"user_agent"`
	IP               string    `json:"ip"`
	ExpiresAt        time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
	LastUsedAt       time.Time `json:"last_used_at"`
}
