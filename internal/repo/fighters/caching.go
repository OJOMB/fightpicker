package fighters

import (
	"time"

	"github.com/gofrs/uuid"
)

const fighterCacheTTL = 5 * time.Minute

func fighterCacheKey(id uuid.UUID) string {
	return "fighter:" + id.String()
}
