package sys_service

import (
	"context"
	"skygo_detection/guardian/app/session"
)

type SessionCacheFun func(args ...interface{}) interface{}

func SessionCacheGet(ctx context.Context, key string, f SessionCacheFun, args ...interface{}) (result interface{}) {
	if val := session.Get(ctx, key); val == nil {
		result = f(args...)
		session.Set(ctx, key, result)
	} else {
		result = val
	}
	return
}
