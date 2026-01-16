package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// PresignedPutURLGenerator defines the interface for generating presigned PUT URLs.
type PresignedPutURLGenerator interface {
	GeneratePresignedPutURL(ctx context.Context, userID id.UUID7, contentType string) (string, http.Header, error)
}

// generatePresignedURL handles the HTTP POST request for the v1 generate_presigned_url endpoint.
func (h *Handler) generatePresignedURL(svc PresignedPutURLGenerator) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		userID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		defer r.Body.Close()
		var userUpdateReq dtos.UserProfilePictureUploadURLRequest
		if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
			return v1.ErrUnreadableRequestBody
		}

		url, headers, err := svc.GeneratePresignedPutURL(ctx, userID, userUpdateReq.ContentType)
		if err != nil {
			return err
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

		h.WriteJSON(ctx, w, http.StatusOK, resp)

		return nil
	}
}
