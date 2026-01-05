package fighters

import (
	"time"

	"github.com/OJOMB/fightpicker/internal/service"
	"github.com/gofrs/uuid"
)

type Fighter struct {
	ID                uuid.UUID
	FirstName         string
	LastName          string
	Nickname          string
	Gender            service.Gender
	DOB               DOB
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
	CreatedBy         uuid.UUID
	UpdatedAt         time.Time
	UpdatedBy         uuid.UUID
}

// DOB is a custom type for handling date of birth with special JSON unmarshalling.
// date in the format "YYYY-MM-DD"
type DOB time.Time

func (d *DOB) UnmarshalJSON(data []byte) error {
	str := string(data)
	t, err := time.Parse(`"2006-01-02"`, str)
	if err != nil {
		return err
	}

	*d = DOB(t)
	return nil
}
