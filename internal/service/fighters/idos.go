package fighters

import (
	"time"

	"github.com/OJOMB/fightpicker/internal/service"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type Fighter struct {
	ID                id.UUID7
	FirstName         string
	LastName          string
	Nickname          string
	Gender            service.Gender
	DOB               time.Time
	Height            float64
	Weight            float64
	Reach             float64
	Stance            string
	Country           string
	FightingOutOf     string
	Bio               string
	ProfilePicture    string
	Wins              int
	Losses            int
	Draws             int
	Disqualifications int
	NoContests        int
	CreatedAt         time.Time
	CreatedBy         id.UUID7
	UpdatedAt         time.Time
	UpdatedBy         id.UUID7
}

type IngestionSummary struct {
	Created int
	Failed  int
	Skipped int
	Updated int
	Total   int
}

type IngestRow struct {
	Index   int
	Fighter Fighter
}

type IngestionResult struct {
	// Index Index of the item in the request array
	Index int

	// FighterId Present when status is created or updated
	FighterId *id.UUID7
	Status    IngestionResultStatus

	Error *ErrorObject `json:"error,omitempty"`
}

type IngestionResultStatus int

const (
	IngestionResultStatusCreated IngestionResultStatus = iota + 1
	IngestionResultStatusUpdated
	IngestionResultStatusSkipped
	IngestionResultStatusFailed
)

func (irs IngestionResultStatus) String() string {
	switch irs {
	case IngestionResultStatusCreated:
		return "created"
	case IngestionResultStatusUpdated:
		return "updated"
	case IngestionResultStatusSkipped:
		return "skipped"
	case IngestionResultStatusFailed:
		return "failed"
	default:
		return ""
	}
}

func IngestionResultStatusFromString(s string) IngestionResultStatus {
	switch s {
	case "created":
		return IngestionResultStatusCreated
	case "updated":
		return IngestionResultStatusUpdated
	case "skipped":
		return IngestionResultStatusSkipped
	case "failed":
		return IngestionResultStatusFailed
	default:
		return 0
	}
}

// Scan implements the sql.Scanner interface for IngestionResultStatus
func (irs *IngestionResultStatus) Scan(value interface{}) error {
	if value == nil {
		*irs = 0
		return nil
	}

	// string to int mapping
	strVal, ok := value.(string)
	if !ok {
		return nil
	}

	intVal := IngestionResultStatusFromString(strVal)

	*irs = IngestionResultStatus(intVal)
	return nil
}

// Value implements the driver.Valuer interface for IngestionResultStatus
func (irs IngestionResultStatus) Value() (interface{}, error) {
	return irs.String(), nil
}

type ErrorObject struct {
	Code    string
	Message string
}
