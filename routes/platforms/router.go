package platforms

import (
	docs "Chlamydia/doclib"
	"Chlamydia/state"
	"Chlamydia/types"
	"Chlamydia/uapi"

	"net/http"

	"github.com/go-chi/chi/v5"
)

const tagName = "Platforms"

type Router struct{}

func (b Router) Tag() (string, string) {
	return tagName, "Endpoints regarding our supported platforms."
}

func Docs() *docs.Doc {
	return &docs.Doc{
		Summary:     "List Platforms",
		Description: "List(s) all of the platforms we support.",
		Resp:        types.ApiError{},
		Params:      []docs.Parameter{},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	type Response struct {
		Platforms []types.PlatformType `json:"platforms"`
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json: &Response{
			Platforms: state.GetPlatforms(),
		},
	}
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/platforms",
		OpId:    "/platforms",
		Method:  uapi.GET,
		Docs:    Docs,
		Handler: Route,
	}.Route(r)
}
