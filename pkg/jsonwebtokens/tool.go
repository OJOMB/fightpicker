package jsonwebtokens

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	TokenStr  string
	JTI       uuid.UUID
	UserID    uuid.UUID
	Issuer    string
	Audience  string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Claims    map[string]any
}

type Claims[T any] struct {
	// unfortunately right now Go doesn't support embedding type parameters so we are forced to name this field "Custom"
	// ideally custom claims would be flattened into the parent struct but this is not currently possible
	// there seems to be proposals to allow this feature in the future though.
	Custom T
	jwt.RegisteredClaims
}

type JWTTool[T any] struct{}

func NewJWTTool[T any]() *JWTTool[T] {
	return &JWTTool[T]{}
}
