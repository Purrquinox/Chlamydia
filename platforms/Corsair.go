package platforms

import (
	"Chlamydia/types"
)

type Corsair struct {
	name string
}

func (c *Corsair) Get(name string) {
	c.name = name
}

func (c *Corsair) SetName(name string) {
	c.name = name
}

func CorsairData() types.PlatformType {
	return types.PlatformType{
		Name:                 "Corsair",
		Description:          "CORSAIR is a leading global developer and manufacturer of high-performance gear and technology for gamers, content creators, and PC enthusiasts.",
		Logo:                 "https://res.cloudinary.com/corsair-pwa/image/upload/f_auto,q_auto/v1665096094/akamai/content/images/reusable/CORSAIR_logo_horizontal_white.png",
		Show:                 true,
		Platform:             &Corsair{},
		AdditionalParameters: []map[string]string{},
	}
}
