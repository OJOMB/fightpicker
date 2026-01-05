package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type PresignedPutURLGenerator interface {
	GeneratePresignedPutURL(ctx context.Context, userID uuid.UUID, contentType string) (string, http.Header, error)
}

// updateUser handles the HTTP PATCH request for the v1 update_user endpoint.
func (h *Handler) generatePresignedURL(svc PresignedPutURLGenerator, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersUpdate)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := mux.Vars(r)[v1.QueryParamUserID]
		userID, err := uuid.FromString(id)
		if err != nil {
			logger.DebugContext(ctx, "invalid user_id parameter", "error", err, "user_id", id)
			h.writeError(ctx, w, errors.Wrap(v1.ErrInvalidUUID, v1.QueryParamUserID), logger)
			return
		}

		var userUpdateReq dtos.UserProfilePictureUploadURLRequest
		if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
			logger.ErrorContext(ctx, "failed to decode request body", "error", err)
		}

		url, headers, err := svc.GeneratePresignedPutURL(ctx, userID, userUpdateReq.ContentType)
		if err != nil {
			logger.ErrorContext(ctx, "failed to generate presigned URL", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		// TODO: not sure if we need to even return headers here, but for now we will
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.ErrorContext(ctx, "failed to encode error response", "error", err)
		}
	}
}
