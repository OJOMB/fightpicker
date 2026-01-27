package fighters

import (
	"time"

	"github.com/oapi-codegen/runtime/types"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	"github.com/OJOMB/fightpicker/internal/service"
	fightersservice "github.com/OJOMB/fightpicker/internal/service/fighters"
)

func fighterIDOToDTO(fighter fightersservice.Fighter) dtos.FighterResponse {
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

func fightersCreateReqDTOToIDO(req dtos.FighterCreateReq) fightersservice.Fighter {
	return fightersservice.Fighter{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Nickname:          req.Nickname,
		Gender:            service.GenderFromString(string(req.Gender)),
		Weight:            req.Weight,
		Height:            req.Height,
		Reach:             req.Reach,
		DOB:               req.Dob.Time,
		Stance:            req.Stance,
		Country:           req.Country,
		FightingOutOf:     req.FightingOutOf,
		Wins:              req.Wins,
		Losses:            req.Losses,
		Draws:             req.Draws,
		Disqualifications: req.Disqualifications,
		NoContests:        req.NoContests,
	}
}

func fightersIngestionReqDTOtoIDO(req dtos.FightersIngestionReq) []fightersservice.IngestRow {
	fighterIngestionRows := make([]fightersservice.IngestRow, len(req))
	for i, fighterDTO := range req {
		fighterIngestionRows[i] = fightersservice.IngestRow{
			Index:   i,
			Fighter: fightersCreateReqDTOToIDO(fighterDTO),
		}
	}

	return fighterIngestionRows
}

func ingestionIDOstoDTO(results []fightersservice.IngestionResult, summary fightersservice.IngestionSummary) dtos.FightersIngestionResp {
	dtoResults := make([]dtos.FighterIngestionResult, len(results))
	for i, result := range results {
		dtoResults[i] = ingestionResultIDOtoDTO(result)
	}

	dtoSummary := ingestionSummaryIDOtoDTO(summary)

	return dtos.FightersIngestionResp{
		Results: dtoResults,
		Summary: dtoSummary,
	}
}

func ingestionResultIDOtoDTO(ido fightersservice.IngestionResult) dtos.FighterIngestionResult {
	dto := dtos.FighterIngestionResult{
		Index:     ido.Index,
		Status:    dtos.FighterIngestionResultStatus(ido.Status.String()),
		FighterId: ido.FighterId,
	}

	if ido.Error != nil {
		dto.Error = &dtos.ErrorObject{
			Code:    ido.Error.Code,
			Message: ido.Error.Message,
		}
	}

	return dto
}

func ingestionSummaryIDOtoDTO(summary fightersservice.IngestionSummary) dtos.FightersIngestionSummary {
	return dtos.FightersIngestionSummary{
		Created: summary.Created,
		Failed:  summary.Failed,
		Skipped: summary.Skipped,
		Updated: summary.Updated,
		Total:   summary.Total,
	}
}
