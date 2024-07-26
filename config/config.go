package config

type Config struct {
	Name        string
	Version     string
	Description string
	Port        string
	Contact     ContactCard
	Links       []Link
}

type ContactCard struct {
	Name  string
	URL   string
	Email string
}

type Link struct {
	Name string
	URL  string
}

func NewConfig() *Config {
	return &Config{
		Name:        "Chlamydia",
		Version:     "0.0.1-predev",
		Description: "Chlamydia is a comprehensive program designed to enhance the customization and personalization of your PC's RGB lighting and AIO cooler. By leveraging the SDKs of supported RGB systems, such as Corsair and NZXT, it provides users with greater control and flexibility over their hardware. This allows for a more tailored and immersive gaming or computing experience, making it easier to achieve the exact look and performance you desire.",
		Port:        "50805",
		Contact: ContactCard{
			Name:  "Purrquinox",
			URL:   "https://purrquinox.com/",
			Email: "support@purrquinox.com",
		},
		Links: []Link{
			{
				Name: "License",
				URL:  "https://github.com/Purrquinox/Chlamydia/blob/production/LICENSE",
			},
		},
	}
}

func (c *Config) GetLinkByName(name string) *Link {
	for _, link := range c.Links {
		if link.Name == name {
			return &link
		}
	}

	return nil
}
