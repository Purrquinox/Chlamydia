package types

type PlatformInterface interface {
	ListDevices() []Device
}

type PlatformType struct {
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	Identifier           string              `json:"identifier"`
	Logo                 string              `json:"logo"`
	Show                 bool                `json:"show"`
	Platform             PlatformInterface   `json:"platform"`
	AdditionalParameters []map[string]string `json:"additional_parameters"`
}

type DeviceType string

const (
	All              DeviceType = "All"
	Unknown          DeviceType = "Unknown"
	Keyboard         DeviceType = "Keyboard"
	Mouse            DeviceType = "Mouse"
	Mousemat         DeviceType = "Mousemat"
	Headset          DeviceType = "Headset"
	HeadsetStand     DeviceType = "HeadsetStand"
	FanLedController DeviceType = "FanLedController"
	LedController    DeviceType = "LedController"
	MemoryModule     DeviceType = "MemoryModule"
	Cooler           DeviceType = "Cooler"
	Motherboard      DeviceType = "Motherboard"
	GraphicsCard     DeviceType = "GraphicsCard"
	Touchbar         DeviceType = "Touchbar"
	GameController   DeviceType = "GameController"
)

type Device struct {
	DeviceName     string     `json:"name"`
	DeviceID       string     `json:"id"`
	DeviceType     DeviceType `json:"type"`
	DevicePlatform string     `json:"platform"`
}
