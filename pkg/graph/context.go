package graph

import "context"

type GraphContextKey[T any] struct {
	Key     string
	Default T
}

var DepthKey = GraphContextKey[int]{
	Key:     "depth",
	Default: 10,
}

func WithDepth(ctx context.Context, depth int) context.Context {
	return context.WithValue(ctx, DepthKey, depth)
}

func DepthValue(ctx context.Context) int {
	if d := ctx.Value(DepthKey); d != nil {
		return d.(int)
	}
	return DepthKey.Default
}

var NamespaceFilterKey = GraphContextKey[string]{
	Key:     "ns-filter",
	Default: "",
}

var LimitKey = GraphContextKey[int64]{
	Key:     "limit",
	Default: 0,
}

var OffsetKey = GraphContextKey[int64]{
	Key:     "offset",
	Default: 0,
}

func WithNamespaceFilter(ctx context.Context, ns string) context.Context {
	return context.WithValue(ctx, NamespaceFilterKey, ns)
}

func WithLimit(ctx context.Context, limit int64) context.Context {
	return context.WithValue(ctx, LimitKey, limit)
}

func WithOffset(ctx context.Context, offset int64) context.Context {
	return context.WithValue(ctx, OffsetKey, offset)
}

func LimitValue(ctx context.Context) int64 {
	if v := ctx.Value(LimitKey); v != nil {
		return v.(int64)
	}
	return LimitKey.Default
}

func OffsetValue(ctx context.Context) int64 {
	if v := ctx.Value(OffsetKey); v != nil {
		return v.(int64)
	}
	return OffsetKey.Default
}

func NSFilterValue(ctx context.Context) string {
	if v := ctx.Value(NamespaceFilterKey); v != nil {
		return v.(string)
	}
	return NamespaceFilterKey.Default
}
