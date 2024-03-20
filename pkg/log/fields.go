package log

import (
	"github.com/infinimesh/proto/node/access"
	"go.uber.org/zap"
)

func ZapAccess(level access.Level, role access.Role) zap.Field {
	return zap.Dict("access", zap.String("level", level.String()), zap.String("role", role.String()))
}
