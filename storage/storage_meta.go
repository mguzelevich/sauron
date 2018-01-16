package storage

import (
	"time"
)

type MetaInfo struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
