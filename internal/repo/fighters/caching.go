package fighters

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const fighterCacheTTL = 5 * time.Minute

func fighterCacheKey(id uuid.UUID) string {
	return "fighter:" + id.String()
}
