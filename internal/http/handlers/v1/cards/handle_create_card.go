package cards

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/cards"
)

type CardCreator interface {
	CreateCard(ctx context.Context, card service.Card) (service.Card, error)
}

// createCard handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) createCard(svc CardCreator) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		defer r.Body.Close()
		var cardCreateReq dtos.CardCreateReq
		if err := json.NewDecoder(r.Body).Decode(&cardCreateReq); err != nil {
			// TODO: return proper error response
			return err
		}

		createdCard, err := svc.CreateCard(ctx, cardCreateRequestToIDO(cardCreateReq))
		if err != nil {
			return err
		}

		h.Write(ctx, w, http.StatusCreated, cardIDOToDTO(createdCard))

		return nil
	}
}
