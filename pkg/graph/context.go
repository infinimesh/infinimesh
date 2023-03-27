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

func WithNamespaceFilter(ctx context.Context, ns string) context.Context {
	return context.WithValue(ctx, NamespaceFilterKey, ns)
}

func NSFilterValue(ctx context.Context) string {
	if v := ctx.Value(NamespaceFilterKey); v != nil {
		return v.(string)
	}
	return NamespaceFilterKey.Default
}
