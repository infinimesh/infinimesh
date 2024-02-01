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

var OffsetKey = GraphContextKey[int]{
	Key:     "offset",
	Default: 10,
}

var LimitKey = GraphContextKey[int]{
	Key:     "limit",
	Default: 10,
}

func WithDepth(ctx context.Context, depth int) context.Context {
	return context.WithValue(ctx, DepthKey, depth)
}

func WithOffset(ctx context.Context, limit int) context.Context {
	return context.WithValue(ctx, OffsetKey, limit)
}

func WithLimit(ctx context.Context, limit int) context.Context {
	return context.WithValue(ctx, LimitKey, limit)
}

func DepthValue(ctx context.Context) int {
	if d := ctx.Value(DepthKey); d != nil {
		return d.(int)
	}
	return DepthKey.Default
}

func OffsetValue(ctx context.Context) int {
	if d := ctx.Value(OffsetKey); d != nil {
		return d.(int)
	}
	return 0
}

func LimitValue(ctx context.Context) int {
	if d := ctx.Value(LimitKey); d != nil {
		return d.(int)
	}
	return 0
}

var NamespaceFilterKey = GraphContextKey[string]{
	Key:     "ns-filter",
	Default: "",
}

func WithNamespaceFilter(ctx context.Context, ns string) context.Context {
	return context.WithValue(ctx, NamespaceFilterKey, ns)
}

func NSFilterValue(ctx context.Context) string {
	if v := ctx.Value(NamespaceFilterKey); v != nil {
		return v.(string)
	}
	return NamespaceFilterKey.Default
}
