package model

import (
	"time"
)

type Otp struct {
	ID        uint `gorm:"primaryKey"`
	Code      int
	Attempts  int
	LastRetry time.Time
}