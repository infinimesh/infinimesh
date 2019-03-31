package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	dlq := "default_or_dlq"
	router := New(dlq)

	require.Equal(t, "shadow.reported-state.delta", router.Route("devices/0x1/state/reported/delta", "0x1"))
	require.Equal(t, dlq, router.Route("devices/0x1/state/reported/delta", "0x2"))
	require.Equal(t, dlq, router.Route("/devices/0x1/state/reported/delta", "0x1"))
	require.Equal(t, dlq, router.Route("", "0x1"))
	require.Equal(t, dlq, router.Route("/", "0x1"))
	require.Equal(t, dlq, router.Route("/abc", "0x1"))
	require.Equal(t, dlq, router.Route("/devices/0x1", "0x1"))
	require.Equal(t, dlq, router.Route("/devices/0x1/", "0x1"))
}
