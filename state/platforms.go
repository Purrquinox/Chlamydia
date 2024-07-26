package state

import (
	"Chlamydia/platforms"
	"Chlamydia/types"
)

func GetPlatforms() []types.PlatformType {
	var platformsSlice []types.PlatformType
	platformsSlice = append(platformsSlice, platforms.NZXTData(), platforms.CorsairData())
	return platformsSlice
}
