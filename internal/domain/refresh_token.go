package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         uuid.UUID  `json:"id"`
	EmployeeID uuid.UUID  `json:"employee_id"`
	TokenHash  string     `json:"-"` // Никогда не выводить
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
}

func NewRefreshToken(employeeID uuid.UUID, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) *RefreshToken {
	return &RefreshToken{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
	}
}

func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}

func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.RevokedAt = &now
}
