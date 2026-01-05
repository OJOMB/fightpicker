package fighters

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"

	service "github.com/OJOMB/fightpicker/internal/service/fighters"
)

// GetFighterByID retrieves a fighter by their ID.
func (r *Repo) GetFighterByID(ctx context.Context, id uuid.UUID) (service.Fighter, error) {
	cacheKey := fighterCacheKey(id)
	if r.cachingEnabled {
		cached, err := r.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			var fighter service.Fighter
			err = json.Unmarshal([]byte(cached), &fighter)
			if err == nil {
				return fighter, nil
			}
		}

		r.logger.DebugContext(ctx, "failed to get fighter from cache", "error", err)
	}

	fighter, err := r.client.GetFighterByID(ctx, id)
	if err != nil {
		return service.Fighter{}, dbErrorToServiceError(err)
	}

	if r.cachingEnabled {
		if bytes, err := json.Marshal(fighter); err == nil {
			err := r.cache.Set(ctx, cacheKey, bytes, fighterCacheTTL).Err()
			if err != nil {
				r.logger.WarnContext(ctx, "failed to set fighter in cache", "error", err)
			}
		}
	}

	return fighterDBOtoFighterIDO(fighter), nil
}
