package users

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/pkg/id"
)

type PresignedPutURLGenerator interface {
	GeneratePresignedPutURL(ctx context.Context, userID id.UUID7, contentType string) (string, http.Header, error)
}

func (s *Service) GeneratePresignedPutURL(ctx context.Context, userID id.UUID7, contentType string) (string, http.Header, error) {
	if reqSubject := s.ctxTool.GetReqSubjectFromContext(ctx); reqSubject != userID {
		return "", nil, ErrUnauthorized
	}

	return s.repo.GeneratePresignedPutURL(ctx, userID, contentType)
}
