package jsonwebtokens

import (
	"time"

	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/golang-jwt/jwt/v5"
)

func (jwtt *JWTTool[T]) GenerateToken(jti, userID id.UUID7, duration time.Duration, iss, aud string, customClaims map[string]any, secretKey []byte) (Token, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(duration)

	tokenClaims := jwt.MapClaims{
		"sub":    userID.String(),
		"exp":    expiresAt.Unix(),
		"iat":    now.Unix(),
		"nbf":    now.Unix(),
		"jti":    jti.String(),
		"iss":    iss,
		"aud":    aud,
		"custom": customClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return Token{}, err
	}

	return Token{
		TokenStr:  tokenStr,
		JTI:       jti,
		UserID:    userID,
		Issuer:    iss,
		Audience:  aud,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
		Claims:    tokenClaims,
	}, nil
}
