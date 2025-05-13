package business

import (
	"context"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/redis_cache"
)

func GetBundle(ctx context.Context, bundleId string) (*kxmj_core.ConfigBundle, error) {
	return redis_cache.GetCache().GetBundleCache().GetDetailCache().Get(ctx, bundleId)
}
