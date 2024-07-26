package types

type PlatformInterface interface {
	Get(name string)
	SetName(name string)
}

type PlatformType struct {
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	Logo                 string              `json:"logo"`
	Show                 bool                `json:"show"`
	Platform             PlatformInterface   `json:"platform"`
	AdditionalParameters []map[string]string `json:"additional_parameters"`
}
