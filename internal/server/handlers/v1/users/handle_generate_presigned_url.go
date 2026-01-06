package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// PresignedPutURLGenerator defines the interface for generating presigned PUT URLs.
type PresignedPutURLGenerator interface {
	GeneratePresignedPutURL(ctx context.Context, userID uuid.UUID, contentType string) (string, http.Header, error)
}

// generatePresignedURL handles the HTTP POST request for the v1 generate_presigned_url endpoint.
func (h *Handler) generatePresignedURL(svc PresignedPutURLGenerator, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersUpdate)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := h.parseUserID(r)
		if err != nil {
			h.writeError(ctx, w, errors.Wrap(v1.ErrInvalidUUID, v1.QueryParamUserID), logger)
			return
		}

		defer r.Body.Close()
		var userUpdateReq dtos.UserProfilePictureUploadURLRequest
		if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
			h.writeError(ctx, w, errors.Wrap(v1.ErrUnreadableRequestBody, v1.QueryParamUserID), logger)
			return
		}

		url, headers, err := svc.GeneratePresignedPutURL(ctx, userID, userUpdateReq.ContentType)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		var signedHeaders = make(map[string]string)
		for key, values := range headers {
			if len(values) > 0 {
				signedHeaders[key] = values[0]
			}
		}

		resp := dtos.UserProfilePictureUploadURLResponse{
			UploadUrl:    url,
			SignedHeader: signedHeaders,
		}

		h.writeJSON(ctx, w, logger, http.StatusOK, resp)
	}
}
