package jsonwebtokens

import (
	"crypto/sha256"
	"encoding/hex"
)

func (jwtt *JWTTool[T]) HashTokenString(tokenStr string) string {
	h := sha256.Sum256([]byte(tokenStr))
	return hex.EncodeToString(h[:])
}
