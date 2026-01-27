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

func (f *Fighter) normalize() {
	f.FirstName = service.NormalizeName(f.FirstName)
	f.LastName = service.NormalizeName(f.LastName)
	f.Nickname = service.NormalizeName(f.Nickname)
	f.Stance = service.NormalizeString(f.Stance)
	f.Country = service.NormalizeString(f.Country)
	f.FightingOutOf = service.NormalizeString(f.FightingOutOf)
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

type IngestionResultStatus string

const (
	IngestionResultStatusCreated IngestionResultStatus = "created"
	IngestionResultStatusUpdated IngestionResultStatus = "updated"
	IngestionResultStatusSkipped IngestionResultStatus = "skipped"
	IngestionResultStatusFailed  IngestionResultStatus = "failed"
)

func (s IngestionResultStatus) String() string {
	return string(s)
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
		return ""
	}
}

type ErrorObject struct {
	Code    string
	Message string
}
