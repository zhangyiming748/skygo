package metadata

import (
	"context"
	"fmt"
	"skygo_detection/guardian/src/net/qmap"
	"strconv"
)

type mdKey struct{}

// Pairs returns an MD formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
func Pairs(kv ...interface{}) qmap.QM {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := qmap.QM{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s.(string)
			continue
		}
		md[key] = s
	}
	return md
}

// NewContext creates a new context with md attached.
func NewContext(ctx context.Context, md qmap.QM) context.Context {
	return context.WithValue(ctx, mdKey{}, md)
}

// FromContext returns the incoming metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromContext(ctx context.Context) (md qmap.QM, ok bool) {
	md, ok = ctx.Value(mdKey{}).(qmap.QM)
	return
}

// String get string value from metadata in context
func String(ctx context.Context, key string) string {
	md, ok := ctx.Value(mdKey{}).(qmap.QM)
	if !ok {
		return ""
	}
	return md.String(key)
}

// Int64 get int64 value from metadata in context
func Int64(ctx context.Context, key string) int64 {
	md, ok := ctx.Value(mdKey{}).(qmap.QM)
	if !ok {
		return 0
	}
	return md.Int64(key)
}

// Value get value from metadata in context return nil if not found
func Value(ctx context.Context, key string) interface{} {
	md, ok := ctx.Value(mdKey{}).(qmap.QM)
	if !ok {
		return nil
	}
	return md.Interface(key)
}

// WithContext return no deadline context and retain metadata.
func WithContext(c context.Context) context.Context {
	md, ok := FromContext(c)
	if ok {
		nmd := md.Copy()
		return NewContext(context.Background(), nmd)
	}
	return context.Background()
}

// Bool get boolean from metadata in context use strconv.Parse.
func Bool(ctx context.Context, key string) bool {
	md, ok := ctx.Value(mdKey{}).(qmap.QM)
	if !ok {
		return false
	}
	switch md[key].(type) {
	case bool:
		return md[key].(bool)
	case string:
		ok, _ = strconv.ParseBool(md[key].(string))
		return ok
	default:
		return false
	}
}

func Set(ctx context.Context, key string, val interface{}) bool {
	md, ok := ctx.Value(mdKey{}).(qmap.QM)
	if ok {
		md[key] = val
		return true
	}
	return false
}
