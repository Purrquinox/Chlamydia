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

// Get all platforms that we support
func GetPlatforms() []types.PlatformType {
	var platformsSlice []types.PlatformType
	platformsSlice = append(platformsSlice, platforms.NZXTData(), platforms.CorsairData(), platforms.TestData())
	return platformsSlice
}

// Get all devices that were discovered. (only shows devices that are with platforms we support)
func GetDevices() []types.Device {
	devices := []types.Device{}
	for _, platform := range GetPlatforms() {
		devices = append(devices, platform.Platform.ListDevices()...)
	}
	return devices
}

// Get all devices within a platform that were discovered.
func GetDevicesByPlatform(id string) []types.Device {
	devices := []types.Device{}
	for _, platform := range GetPlatforms() {
		devices = append(devices, platform.Platform.ListDevices()...)
	}
	return devices
}

func Setup() {
	Logger = snippets.CreateZap()
}
