package apierr

import (
	"github.com/OJOMB/fightpicker/internal/http/dtos"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (err *APIError) ToDTO(reqID id.UUID7) dtos.ErrorEnvelope {
	return dtos.ErrorEnvelope{
		Error: dtos.ErrorObject{
			RequestId: reqID.String(),
			Code:      err.Code,
			Message:   err.Public.Error(),
		},
	}
}
