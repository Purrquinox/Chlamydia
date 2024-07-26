package state

import (
	"context"

	"github.com/infinitybotlist/eureka/snippets"
	"go.uber.org/zap"
)

var (
	Logger  *zap.Logger
	Context = context.Background()
)

func Setup() {
	Logger = snippets.CreateZap()
}
