package state

import (
	"Chlamydia/platforms"
	"Chlamydia/types"
	"context"

	"github.com/infinitybotlist/eureka/snippets"
	"go.uber.org/zap"
)

var (
	Logger  *zap.Logger
	Context = context.Background()
)

func GetPlatforms() []types.PlatformType {
	var platformsSlice []types.PlatformType
	platformsSlice = append(platformsSlice, platforms.NZXTData(), platforms.CorsairData())
	return platformsSlice
}

func Setup() {
	Logger = snippets.CreateZap()
}
