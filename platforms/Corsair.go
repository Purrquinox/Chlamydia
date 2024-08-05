package platforms

import (
	"Chlamydia/types"
	"github.com/google/uuid"
)

type Corsair struct{}

var CorsairID = uuid.New().String()

func (c *Corsair) ListDevices() []types.Device {
	devices := []types.Device{
		{
			DeviceName:     "Vengence Pro",
			DeviceID:       "vengence_pro",
			DeviceType:     types.MemoryModule,
			DevicePlatform: CorsairID,
		},
		{
			DeviceName:     "Commander Core",
			DeviceID:       "commander_core",
			DeviceType:     types.LedController,
			DevicePlatform: CorsairID,
		},
	}
	return devices
}

func CorsairData() types.PlatformType {
	return types.PlatformType{
		Name:                 "Corsair",
		Description:          "CORSAIR is a leading global developer and manufacturer of high-performance gear and technology for gamers, content creators, and PC enthusiasts.",
		Identifier:           CorsairID,
		Logo:                 "https://res.cloudinary.com/corsair-pwa/image/upload/f_auto,q_auto/v1665096094/akamai/content/images/reusable/CORSAIR_logo_horizontal_white.png",
		Show:                 true,
		Platform:             &Corsair{},
		AdditionalParameters: []map[string]string{},
	}
}
