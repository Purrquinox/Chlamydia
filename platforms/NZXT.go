package platforms

import (
	"Chlamydia/types"
)

type NZXT struct {
	name string
}

func (c *NZXT) Get(name string) {
	c.name = name
}

func (c *NZXT) SetName(name string) {
	c.name = name
}

func NZXTData() types.PlatformType {
	return types.PlatformType{
		Name:                 "NZXT",
		Description:          "NZXT is a leading force in PC hardware, known for its groundbreaking designs and high-performance components. Their products, ranging from sleek cases to advanced cooling solutions, are crafted to elevate both the aesthetics and functionality of your gaming and workstation setups.",
		Logo:                 "https://nzxt.com/assets/cms/34299/1611033291-nzxt-logo.png?auto=format&fit=max&h=540&w=540",
		Show:                 true,
		Platform:             &NZXT{},
		AdditionalParameters: []map[string]string{},
	}
}
