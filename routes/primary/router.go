package primary

import (
	"Chlamydia/config"
	docs "Chlamydia/doclib"
	"Chlamydia/state"
	"Chlamydia/types"
	"Chlamydia/uapi"

	"net/http"

	"github.com/go-chi/chi/v5"
)

const tagName = "Main"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "Main page of Chlamydia Core (/)"
}

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Main page",
		Description: "Lists basic information about the Chlamydia Core",
		Resp:        types.ApiError{},
		Params:      []docs.Parameter{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	type Response struct {
		Name              string               `json:"name"`
		Version           string               `json:"version"`
		Description       string               `json:"description"`
		Port              string               `json:"port"`
		OpenAPI           string               `json:"openapi"`
		DocumentationLink string               `json:"documentation_link"`
		Platforms         []types.PlatformType `json:"platforms"`
	}

	config := config.NewConfig()
	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json: &Response{
			Name:              config.Name,
			Version:           config.Version,
			Description:       config.Description,
			Port:              config.Port,
			OpenAPI:           "http://localhost:" + config.Port + "/openapi",
			DocumentationLink: "http://localhost:" + config.Port + "/docs",
			Platforms:         state.GetPlatforms(),
		},
	}
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/",
		OpId:    "/",
		Method:  uapi.GET,
		Docs:    Docs,
		Handler: Route,
	}.Route(r)
}
