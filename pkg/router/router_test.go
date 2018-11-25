package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	routes := map[string]string{
		"/shadows/":             "shadows",
		"/shadows/ab/":          "shadow-ab",
		"/shadows/abc/":         "shadow-abc",
		"/shadows/abc/specific": "more specific route",
		"/devices/":             "devices",
	}
	router := New("default_or_dlq", routes)

	require.Equal(t, "shadow-abc", router.Route("/shadows/abc/xx"))
	require.Equal(t, "more specific route", router.Route("/shadows/abc/specific"))
	require.Equal(t, "default_or_dlq", router.Route("whatever"))
}
