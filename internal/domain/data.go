package domain

import (
	"time"

	"github.com/google/uuid"
)

type FileData struct {
	ID        uuid.UUID
	Name      string
	Active    bool
	Temp      float64
	Tags      []string
	Metadata  map[string]interface{}
	CreatedAt time.Time
}
