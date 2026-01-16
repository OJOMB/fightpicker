package fighters

import (
	"time"

	"github.com/OJOMB/fightpicker/pkg/id"
)

const fighterCacheTTL = 5 * time.Minute

func fighterCacheKey(id id.UUID7) string {
	return "fighter:" + id.String()
}
