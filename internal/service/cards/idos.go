package cards

import (
	"time"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// Card represents a fight card with metadata.
type Card struct {
	ID        id.UUID7
	Name      string
	Location  string
	Date      time.Time
	Promotion string

	CreatedAt time.Time
	CreatedBy id.UUID7
	UpdatedAt time.Time
	UpdatedBy id.UUID7
}
