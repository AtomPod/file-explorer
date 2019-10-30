package models

import "time"

//VerificationCode 验证码
type VerificationCode struct {
	ID         int           `json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	Type       string        `json:"type"`
	Code       string        `json:"code"`
	Target     string        `json:"target"`
	Expiration time.Duration `json:"expiration "`
}
