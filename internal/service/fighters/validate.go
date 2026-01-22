package fighters

import (
	"github.com/pkg/errors"
)

// type Fighter struct {
// 	ID                id.UUID7
// 	FirstName         string
// 	LastName          string
// 	Nickname          string
// 	Gender            service.Gender
// 	DOB               time.Time
// 	Height            float64
// 	Weight            float64
// 	Reach             float64
// 	Stance            string
// 	Country           string
// 	FightingOutOf     string
// 	Bio               string
// 	ProfilePicture    string
// 	Wins              int
// 	Losses            int
// 	Draws             int
// 	Disqualifications int
// 	NoContests        int
// 	CreatedAt         time.Time
// 	CreatedBy         id.UUID7
// 	UpdatedAt         time.Time
// 	UpdatedBy         id.UUID7
// }

func (s *Service) validateCreationReq(f *Fighter) error {
	if f.FirstName == "" {
		return errors.Wrap(ErrMissingParameter, "first_name")
	}

	if f.LastName == "" {
		return errors.Wrap(ErrMissingParameter, "last_name")
	}

	if f.DOB.IsZero() {
		return errors.Wrap(ErrMissingParameter, "dob")
	}

	if f.Height <= 0 {
		return errors.Wrap(ErrMissingParameter, "height")
	}

	if f.Weight <= 0 {
		return errors.Wrap(ErrMissingParameter, "weight")
	}

	if f.Reach <= 0 {
		return errors.Wrap(ErrMissingParameter, "reach")
	}

	if f.Stance == "" {
		return errors.Wrap(ErrMissingParameter, "stance")
	}

	if f.Country == "" {
		return errors.Wrap(ErrMissingParameter, "country")
	}

	if f.FightingOutOf == "" {
		return errors.Wrap(ErrMissingParameter, "fighting_out_of")
	}

	if f.Wins < 0 {
		return errors.Wrap(ErrInvalidParameter, "wins")
	}

	if f.Losses < 0 {
		return errors.Wrap(ErrInvalidParameter, "losses")
	}

	if f.Draws < 0 {
		return errors.Wrap(ErrInvalidParameter, "draws")
	}

	if f.Disqualifications < 0 {
		return errors.Wrap(ErrInvalidParameter, "disqualifications")
	}

	if f.NoContests < 0 {
		return errors.Wrap(ErrInvalidParameter, "no_contests")
	}

	return nil
}
