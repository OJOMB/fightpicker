package cards

import (
	"github.com/OJOMB/fightpicker/internal/http/dtos"
	cardservice "github.com/OJOMB/fightpicker/internal/service/cards"
)

func cardCreateRequestToIDO(dto dtos.CardCreateReq) cardservice.Card {
	return cardservice.Card{
		Name:      dto.Name,
		Location:  dto.Location,
		Date:      dto.Date,
		Promotion: dto.Promotion,
	}
}

func cardIDOToDTO(c cardservice.Card) dtos.CardResponse {
	return dtos.CardResponse{
		Id:        c.ID,
		Name:      c.Name,
		Location:  c.Location,
		Date:      c.Date,
		Promotion: c.Promotion,

		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy,
	}
}
