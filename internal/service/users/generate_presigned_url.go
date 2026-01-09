package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

type PresignedPutURLGenerator interface {
	GeneratePresignedPutURL(ctx context.Context, userID uuid.UUID, contentType string) (string, http.Header, error)
}

func (s *Service) GeneratePresignedPutURL(ctx context.Context, userID uuid.UUID, contentType string) (string, http.Header, error) {
	if reqSubject := contextual.GetReqSubjectFromContext(ctx); reqSubject != userID {
		return "", nil, ErrUnauthorized
	}

	return s.repo.GeneratePresignedPutURL(ctx, userID, contentType)
}
