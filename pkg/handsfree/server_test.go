package handsfree

import (
	"testing"

	"go.uber.org/zap"
)

func Test_JustTest(t *testing.T) {
	srv := NewHandsfreeServer(zap.NewExample())
	code := GenerateCode(srv.db)
	if len(code) != 6 {
		t.Fatalf("Code format is wrong: %s", code)
	}
	t.Logf("Just random Code: %s", code)
}
