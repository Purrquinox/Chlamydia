package platforms

import (
	"Chlamydia/types"
	"github.com/google/uuid"
)

type NZXT struct{}

var NZXTID = uuid.New().String()

func (c *NZXT) ListDevices() []types.Device {
	devices := []types.Device{
		{
			DeviceName:     "NZXT Kraken 240 LCD",
			DeviceID:       "kraken_240",
			DeviceType:     types.Unknown,
			DevicePlatform: NZXTID,
		},
	}
	return devices
}

func NZXTData() types.PlatformType {
	return types.PlatformType{
		Name:                 "NZXT",
		Description:          "NZXT is a leading force in PC hardware, known for its groundbreaking designs and high-performance components. Their products, ranging from sleek cases to advanced cooling solutions, are crafted to elevate both the aesthetics and functionality of your gaming and workstation setups.",
		Identifier:           NZXTID,
		Logo:                 "https://nzxt.com/assets/cms/34299/1611033291-nzxt-logo.png?auto=format&fit=max&h=540&w=540",
		Show:                 true,
		Platform:             &NZXT{},
		AdditionalParameters: []map[string]string{},
	}
}
