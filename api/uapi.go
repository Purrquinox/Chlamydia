package api

import (
	"Chlamydia/constants"
	"Chlamydia/state"
	"Chlamydia/types"
	"net/http"

	"Chlamydia/uapi"
)

type DefaultResponder struct{}

func (d DefaultResponder) New(err string, ctx map[string]string) any {
	return types.ApiError{
		Message: err,
		Context: ctx,
	}
}

func Setup() {
	uapi.SetupState(uapi.UAPIState{
		Logger:  state.Logger,
		Context: state.Context,
		Constants: &uapi.UAPIConstants{
			ResourceNotFound:    constants.ResourceNotFound,
			BadRequest:          constants.BadRequest,
			Forbidden:           constants.Forbidden,
			Unauthorized:        constants.Unauthorized,
			InternalServerError: constants.InternalServerError,
			MethodNotAllowed:    constants.MethodNotAllowed,
			BodyRequired:        constants.BodyRequired,
		},
		DefaultResponder: DefaultResponder{},
		Authorize: func(r uapi.Route, req *http.Request) (uapi.AuthData, uapi.HttpResponse, bool) {
			return uapi.AuthData{}, uapi.HttpResponse{}, true
		},
	})
}
