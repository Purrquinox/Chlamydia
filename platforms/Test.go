package platforms

import (
	"Chlamydia/types"
	"github.com/google/uuid"
)

type Test struct{}

var TestID = uuid.New().String()

func (c *Test) ListDevices() []types.Device {
	devices := []types.Device{
		{
			DeviceName:     "Test Device 1",
			DeviceID:       "test_device_1",
			DeviceType:     types.Motherboard,
			DevicePlatform: TestID,
		},
		{
			DeviceName:     "Test Device 2",
			DeviceID:       "test_device_2",
			DeviceType:     types.LedController,
			DevicePlatform: TestID,
		},
		{
			DeviceName:     "Test Device 3",
			DeviceID:       "test_device_3",
			DeviceType:     types.Cooler,
			DevicePlatform: TestID,
		},
	}
	return devices
}

func TestData() types.PlatformType {
	return types.PlatformType{
		Name:                 "Test Platform",
		Description:          "This is a Test Platform that ensures that everything is working properly. This platform is primarily used during development.",
		Identifier:           TestID,
		Logo:                 "https://res.cloudinary.com/Test-pwa/image/upload/f_auto,q_auto/v1665096094/akamai/content/images/reusable/Test_logo_horizontal_white.png",
		Show:                 true,
		Platform:             &Test{},
		AdditionalParameters: []map[string]string{},
	}
}
