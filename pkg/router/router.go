package router

import (
	radix "github.com/armon/go-radix"
)

type Router struct {
	radix *radix.Tree
}

func New(fallback string, routes map[string]string) *Router {
	routez := make(map[string]interface{})
	for key, value := range routes {
		routez[key] = value
	}
	tree := radix.NewFromMap(routez)
	tree.Insert("", fallback)

	return &Router{
		radix: tree,
	}
}

// TODO force trailing backslash? otherwise, /devices/abc/ may be matched by
// /devices/ab, even if's a completely different device. so we would have to
// honor/enforce the separator. while this may not be an error, it may not be
// what the caller expects.
func (r *Router) Route(inputTopic string) (outputTopic string) {
	_, target, _ := r.radix.LongestPrefix(inputTopic)
	return target.(string)
}
