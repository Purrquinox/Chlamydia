package getPlatform

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
		Summary:     "Get Platform",
		Description: "Get one of the platforms we support.",
		Resp:        types.ApiError{},
		Params: []docs.Parameter{
			{
				Name:        "platform_name",
				Description: "Name of the platform to get.",
				In:          "querystring",
				Required:    true,
				Schema:      "?platform_name=0",
			},
		},
	}
}

func Route(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	query := r.URL.Query().Get("platform_name")

	type Response struct {
		Platform types.PlatformType `json:"platform"`
		Devices  []types.Device     `json:"devices"`
	}

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json: &Response{
			Platform: state.GetPlatform(query),
			Devices:  state.GetDevicesByPlatform(query),
		},
	}
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/platforms/get",
		OpId:    "/platforms/get",
		Method:  uapi.GET,
		Docs:    Docs,
		Handler: Route,
	}.Route(r)
}
