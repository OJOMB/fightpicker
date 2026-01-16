package auth

import (
	"encoding/json"
	"time"

	"github.com/OJOMB/fightpicker/pkg/id"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// Permissions represents a nested map of user permissions organized by version, resource, operation to a set of named permissions.
type Permissions map[string]map[string]map[string]map[string]struct{}

// UnmarshalJSON implements the json.Unmarshaler interface for Permissions.
func (p *Permissions) UnmarshalJSON(data []byte) error {
	raw := map[string]map[string]map[string]map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert `any` → `struct{}`
	perms := Permissions{}
	for v, resources := range raw {
		perms[v] = map[string]map[string]map[string]struct{}{}
		for resource, ops := range resources {
			perms[v][resource] = map[string]map[string]struct{}{}
			for op, permNames := range ops {
				perms[v][resource][op] = map[string]struct{}{}
				for permName := range permNames {
					perms[v][resource][op][permName] = struct{}{}
				}
			}
		}
	}

	*p = perms
	return nil
}

// AuthClaims represents the custom claims stored in a JWT token for authentication.
type AuthClaims struct {
	Perms Permissions `json:"perms"`
	Roles []string    `json:"roles"`
}

// RefreshToken represents a stored refresh token.
type RefreshToken struct {
	ID         id.UUID7  `json:"id"`
	UserID     id.UUID7  `json:"user_id"`
	TokenHash  string    `json:"token_hash"`
	JTI        id.UUID7  `json:"jti"`
	Revoked    bool      `json:"revoked"`
	ExpiresAt  time.Time `json:"expires_at"`
	ReplacedBy *id.UUID7 `json:"replaced_by"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
