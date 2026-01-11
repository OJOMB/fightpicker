package fighters

import (
	"time"

	"github.com/oapi-codegen/runtime/types"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	service "github.com/OJOMB/fightpicker/internal/service/fighters"
)

func fighterIDOToDTO(fighter service.Fighter) dtos.FighterResponse {
	return dtos.FighterResponse{
		Id:        fighter.ID,
		FirstName: fighter.FirstName,
		LastName:  fighter.LastName,
		Nickname:  fighter.Nickname,
		Gender:    dtos.FighterResponseGender(fighter.Gender.String()),
		Weight:    fighter.Weight,
		Height:    fighter.Height,
		Reach:     fighter.Reach,
		Dob: types.Date{
			Time: time.Time(fighter.DOB),
		},
		Stance:            fighter.Stance,
		Country:           fighter.Country,
		FightingOutOf:     fighter.FightingOutOf,
		Bio:               fighter.Bio,
		ProfilePicture:    fighter.ProfilePicture,
		Wins:              fighter.Wins,
		Losses:            fighter.Losses,
		Draws:             fighter.Draws,
		Disqualifications: fighter.Disqualifications,
		NoContests:        fighter.NoContests,
		CreatedAt:         fighter.CreatedAt,
		CreatedBy:         fighter.CreatedBy,
		UpdatedAt:         fighter.UpdatedAt,
		UpdatedBy:         fighter.UpdatedBy,
	}
}
