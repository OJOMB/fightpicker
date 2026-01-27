package fighters

import (
	"github.com/OJOMB/fightpicker/internal/service"
	fightersservice "github.com/OJOMB/fightpicker/internal/service/fighters"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

func fighterDBOtoFighterIDO(dbo postgres.Fighter) fightersservice.Fighter {
	var weight float64
	weightF8, err := dbo.Weight.Float64Value()
	if err == nil {
		weight = weightF8.Float64
	}

	var height float64
	heightF8, err := dbo.Height.Float64Value()
	if err == nil {
		height = heightF8.Float64
	}

	var reach float64
	reachF8, err := dbo.Reach.Float64Value()
	if err == nil {
		reach = reachF8.Float64
	}

	return fightersservice.Fighter{
		ID:                dbo.ID,
		FirstName:         dbo.FirstName,
		LastName:          dbo.LastName,
		Nickname:          dbo.Nickname.String,
		Gender:            service.GenderFromString(string(dbo.Gender)),
		DOB:               dbo.Dob.Time,
		Weight:            weight,
		Height:            height,
		Reach:             reach,
		Stance:            dbo.Stance,
		Country:           dbo.Country,
		FightingOutOf:     dbo.FightingOutOf,
		ProfilePicture:    dbo.ProfilePicture.String,
		Wins:              int(dbo.Wins),
		Losses:            int(dbo.Losses),
		Draws:             int(dbo.Draws),
		NoContests:        int(dbo.NoContests),
		Disqualifications: int(dbo.Disqualifications),
		CreatedAt:         dbo.CreatedAt.Time,
		CreatedBy:         dbo.CreatedBy.Bytes,
		UpdatedAt:         dbo.UpdatedAt.Time,
		UpdatedBy:         dbo.UpdatedBy.Bytes,
	}
}
