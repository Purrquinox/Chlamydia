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
	platformList := []types.PlatformType{
		platforms.CorsairData(),
		platforms.NZXTData(),
		platforms.TestData(),
	}
	return platformList
}

// Get platform that we support
func GetPlatform(id string) types.PlatformType {
	var p types.PlatformType
	for _, platform := range GetPlatforms() {
		if platform.Name == id {
			p = platform
		}
	}
	return p
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
	kp := GetPlatform(id)
	devices := []types.Device{}
	for _, device := range GetDevices() {
		if device.DevicePlatform == kp.Identifier {
			devices = append(devices, device)
		}
	}
	return devices
}

func Setup() {
	Logger = snippets.CreateZap()
}
