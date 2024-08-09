package uapi

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	docs "Chlamydia/doclib"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/eureka/jsonimpl"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type UAPIConstants struct {
	ResourceNotFound    string
	BadRequest          string
	Forbidden           string
	Unauthorized        string
	InternalServerError string
	MethodNotAllowed    string
	BodyRequired        string
}

type UAPIDefaultResponder interface {
	New(msg string, ctx map[string]string) any
}

type UAPIInitData struct {
	Tag string
}

type UAPIState struct {
	Logger              *zap.Logger
	Authorize           func(r Route, req *http.Request) (AuthData, HttpResponse, bool)
	AuthTypeMap         map[string]string
	RouteDataMiddleware func(rd *RouteData, req *http.Request) (*RouteData, error)
	BaseSanityCheck     func(r Route) error
	PatchDocs           func(d *docs.Doc) *docs.Doc
	Context             context.Context
	Constants           *UAPIConstants
	DefaultResponder    UAPIDefaultResponder
	InitData            UAPIInitData
}

type APIRouter interface {
	Routes(r *chi.Mux)
	Tag() (string, string)
}

func (s *UAPIState) SetCurrentTag(tag string) {
	s.InitData.Tag = tag
}

func SetupState(s UAPIState) {
	if s.Constants == nil {
		panic("Constants is nil")
	}

	State = &s
}

var (
	State *UAPIState
)

// WebSocket group management
type WebSocketGroup struct {
	Clients map[*websocket.Conn]bool
	Mu      sync.Mutex
}

var wsGroups = make(map[string]*WebSocketGroup)
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func addClientToGroup(groupName string, conn *websocket.Conn) {
	group, exists := wsGroups[groupName]
	if !exists {
		group = &WebSocketGroup{Clients: make(map[*websocket.Conn]bool)}
		wsGroups[groupName] = group
	}

	group.Mu.Lock()
	group.Clients[conn] = true
	group.Mu.Unlock()
}

func removeClientFromGroup(groupName string, conn *websocket.Conn) {
	group, exists := wsGroups[groupName]
	if !exists {
		return
	}

	group.Mu.Lock()
	delete(group.Clients, conn)
	group.Mu.Unlock()
}

func broadcastMessage(groupName string, message []byte) {
	group, exists := wsGroups[groupName]
	if !exists {
		return
	}

	group.Mu.Lock()
	defer group.Mu.Unlock()

	for client := range group.Clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.Close()
			delete(group.Clients, client)
		}
	}
}

type WebSocketRoute struct {
	Route
	GroupName string
	AuthType  string
}

func (wsr WebSocketRoute) WSRoute(ro Router) {
	if wsr.OpId == "" {
		panic("OpId is empty: " + wsr.String())
	}

	if wsr.Handler == nil {
		panic("Handler is nil: " + wsr.String())
	}

	ro.Get(wsr.Pattern, func(w http.ResponseWriter, req *http.Request) {
		handleWebSocket(wsr, w, req)
	})
}

func handleWebSocket(wsr WebSocketRoute, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	resp := make(chan HttpResponse)

	go func() {
		defer func() {
			err := recover()

			if err != nil {
				State.Logger.Error("[uapi/handleWebSocket] WebSocket handler panic'd", zap.String("operationId", wsr.OpId), zap.String("method", req.Method), zap.String("endpointPattern", wsr.Pattern), zap.String("path", req.URL.Path), zap.Any("error", err))
				resp <- HttpResponse{
					Status: http.StatusInternalServerError,
					Data:   State.Constants.InternalServerError,
				}
			}
		}()

		conn, err := wsUpgrader.Upgrade(w, req, nil)
		if err != nil {
			State.Logger.Error("[uapi/handleWebSocket] Failed to upgrade to WebSocket", zap.Error(err))
			resp <- HttpResponse{
				Status: http.StatusInternalServerError,
				Data:   State.Constants.InternalServerError,
			}
			return
		}
		defer conn.Close()

		groupName := wsr.GroupName
		addClientToGroup(groupName, conn)
		defer removeClientFromGroup(groupName, conn)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			broadcastMessage(groupName, message)
		}
	}()

	respond(ctx, w, resp)
}

type Method int

const (
	GET Method = iota
	POST
	PATCH
	PUT
	DELETE
	HEAD
)

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PATCH:
		return "PATCH"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	}

	panic("Invalid method")
}

type AuthType struct {
	URLVar       string
	Type         string
	AllowedScope string
}

type AuthData struct {
	TargetType string         `json:"target_type"`
	ID         string         `json:"id"`
	Authorized bool           `json:"authorized"`
	Banned     bool           `json:"banned"`
	Data       map[string]any `json:"data"`
}

type Route struct {
	Method                Method
	Pattern               string
	Aliases               map[string]string
	OpId                  string
	Handler               func(d RouteData, r *http.Request) HttpResponse
	Setup                 func()
	Docs                  func() *docs.Doc
	Auth                  []AuthType
	ExtData               map[string]any
	AuthOptional          bool
	SanityCheck           func() error
	DisablePathSlashCheck bool
}

type RouteData struct {
	Context context.Context
	Auth    AuthData
	Props   map[string]string
}

type Router interface {
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Head(pattern string, h http.HandlerFunc)
}

func (r Route) String() string {
	return r.Method.String() + " " + r.Pattern + " (" + r.OpId + ")"
}

func (r Route) Route(ro Router) {
	if r.OpId == "" {
		panic("OpId is empty: " + r.String())
	}

	if r.Handler == nil {
		panic("Handler is nil: " + r.String())
	}

	if r.Docs == nil {
		panic("Docs is nil: " + r.String())
	}

	if r.Pattern == "" {
		panic("Pattern is empty: " + r.String())
	}

	if State.InitData.Tag == "" {
		panic("CurrentTag is empty: " + r.String())
	}

	if r.Setup != nil {
		r.Setup()
	}

	if State.BaseSanityCheck != nil {
		err := State.BaseSanityCheck(r)

		if err != nil {
			panic("Base sanity check failed: " + err.Error())
		}
	}

	if r.SanityCheck != nil {
		err := r.SanityCheck()

		if err != nil {
			panic("Sanity check failed: " + r.String())
		}
	}

	docsObj := r.Docs()

	docsObj.Pattern = r.Pattern
	docsObj.OpId = r.OpId
	docsObj.Method = r.Method.String()
	docsObj.Tags = []string{State.InitData.Tag}
	docsObj.AuthType = []string{}

	for _, auth := range r.Auth {
		t, ok := State.AuthTypeMap[auth.Type]

		if !ok {
			panic("Invalid auth type: " + auth.Type)
		}

		docsObj.AuthType = append(docsObj.AuthType, t)
	}

	if State.PatchDocs != nil {
		docsObj = State.PatchDocs(docsObj)
	}

	brStart := strings.Count(r.Pattern, "{")
	brEnd := strings.Count(r.Pattern, "}")
	pathParams := []string{}
	patternParams := []string{}

	for _, param := range docsObj.Params {
		if param.In == "" || param.Name == "" || param.Schema == nil {
			panic("Param is missing required fields: " + r.String())
		}

		if param.In == "path" {
			pathParams = append(pathParams, param.Name)
		}
	}

	if !r.DisablePathSlashCheck {
		for _, param := range strings.Split(r.Pattern, "/") {
			if strings.HasPrefix(param, "{") && strings.HasSuffix(param, "}") {
				patternParams = append(patternParams, param[1:len(param)-1])
			} else if strings.Contains(param, "{") || strings.Contains(param, "}") {
				panic("{ and } in pattern but does not start with it " + r.String())
			}
		}
	}

	if brStart != brEnd {
		panic("Mismatched { and } in pattern: " + r.String())
	}

	if brStart != len(pathParams) {
		panic("Mismatched number of params and { in pattern: " + r.String())
	}

	if !r.DisablePathSlashCheck {
		if !slices.Equal(patternParams, pathParams) {
			panic("Mismatched params in pattern and docs: " + r.String())
		}
	}

	if len(r.Aliases) > 0 {
		docsObj.Description += "\n\nAliases for this endpoint:"
		for pattern, reason := range r.Aliases {
			docsObj.Description += "\n\n" + pattern + " (" + reason + ")"
		}
	}

	docs.Route(docsObj)

	createRouteHandler(r, ro, r.Pattern)

	if len(r.Aliases) > 0 {
		for pattern := range r.Aliases {
			createRouteHandler(r, ro, pattern)
		}
	}
}

func createRouteHandler(r Route, ro Router, pat string) {
	switch r.Method {
	case GET:
		ro.Get(pat, func(w http.ResponseWriter, req *http.Request) {
			handle(r, w, req)
		})
	case POST:
		ro.Post(pat, func(w http.ResponseWriter, req *http.Request) {
			handle(r, w, req)
		})
	case PATCH:
		ro.Patch(pat, func(w http.ResponseWriter, req *http.Request) {
			handle(r, w, req)
		})
	case PUT:
		ro.Put(pat, func(w http.ResponseWriter, req *http.Request) {
			handle(r, w, req)
		})
	case DELETE:
		ro.Delete(pat, func(w http.ResponseWriter, req *http.Request) {
			handle(r, w, req)
		})
	case HEAD:
		ro.Head(pat, func(w http.ResponseWriter, req *http.Request) {
			handle(r, w, req)
		})
	default:
		panic("Unknown method for route: " + r.String())
	}
}

func respond(ctx context.Context, w http.ResponseWriter, data chan HttpResponse) {
	select {
	case <-ctx.Done():
		return
	case msg, ok := <-data:
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(State.Constants.InternalServerError))
		}

		if msg.Redirect != "" {
			msg.Headers = map[string]string{
				"Location":     msg.Redirect,
				"Content-Type": "text/html; charset=utf-8",
			}
			msg.Data = "<a href=\"" + msg.Redirect + "\">Found</a>.\n"
			msg.Status = http.StatusFound
		}

		if len(msg.Headers) > 0 {
			for k, v := range msg.Headers {
				w.Header().Set(k, v)
			}
		}

		if msg.Json != nil {
			bytes, err := jsonimpl.Marshal(msg.Json)

			if err != nil {
				State.Logger.Error("[uapi.respond] Failed to unmarshal JSON response", zap.Error(err), zap.Int("size", len(msg.Data)))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(State.Constants.InternalServerError))
				return
			}

			msg.Json = nil
			msg.Bytes = bytes
		}

		if msg.Status == 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(msg.Status)
		}

		if len(msg.Bytes) > 0 {
			w.Write(msg.Bytes)
		}

		w.Write([]byte(msg.Data))
		return
	}
}

type HttpResponse struct {
	Data     string
	Bytes    []byte
	Json     any
	Headers  map[string]string
	Status   int
	Redirect string
}

func CompileValidationErrors(payload any) map[string]string {
	var errors = make(map[string]string)

	structType := reflect.TypeOf(payload)

	for _, f := range reflect.VisibleFields(structType) {
		errors[f.Name] = f.Tag.Get("msg")

		arrayMsg := f.Tag.Get("amsg")

		if arrayMsg != "" {
			errors[f.Name+"$arr"] = arrayMsg
		}
	}

	return errors
}

func ValidatorErrorResponse(compiled map[string]string, v validator.ValidationErrors) HttpResponse {
	var errors = make(map[string]string)

	firstError := ""

	for i, err := range v {
		fname := err.StructField()
		if strings.Contains(err.Field(), "[") {
			fname = strings.Split(err.Field(), "[")[0] + "$arr"
		}

		field := compiled[fname]

		var errorMsg string
		if field != "" {
			errorMsg = field + " [" + err.Tag() + "]"
		} else {
			errorMsg = err.Error()
		}

		if i == 0 {
			firstError = errorMsg
		}

		errors[err.StructField()] = errorMsg
	}

	return HttpResponse{
		Status: http.StatusBadRequest,
		Json:   State.DefaultResponder.New(firstError, errors),
	}
}

func DefaultResponse(statusCode int) HttpResponse {
	switch statusCode {
	case http.StatusForbidden:
		return HttpResponse{
			Status: statusCode,
			Data:   State.Constants.Forbidden,
		}
	case http.StatusUnauthorized:
		return HttpResponse{
			Status: statusCode,
			Data:   State.Constants.Unauthorized,
		}
	case http.StatusNotFound:
		return HttpResponse{
			Status: statusCode,
			Data:   State.Constants.ResourceNotFound,
		}
	case http.StatusBadRequest:
		return HttpResponse{
			Status: statusCode,
			Data:   State.Constants.BadRequest,
		}
	case http.StatusInternalServerError:
		return HttpResponse{
			Status: statusCode,
			Data:   State.Constants.InternalServerError,
		}
	case http.StatusMethodNotAllowed:
		return HttpResponse{
			Status: statusCode,
			Data:   State.Constants.MethodNotAllowed,
		}
	case http.StatusNoContent, http.StatusOK:
		return HttpResponse{
			Status: http.StatusNoContent,
		}
	}

	return HttpResponse{
		Status: statusCode,
		Data:   State.Constants.InternalServerError,
	}
}

func handle(r Route, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	resp := make(chan HttpResponse)

	go func() {
		defer func() {
			err := recover()

			if err != nil {
				State.Logger.Error("[uapi/handle] Request handler panic'd", zap.String("operationId", r.OpId), zap.String("method", req.Method), zap.String("endpointPattern", r.Pattern), zap.String("path", req.URL.Path), zap.Any("error", err))
				resp <- HttpResponse{
					Status: http.StatusInternalServerError,
					Data:   State.Constants.InternalServerError,
				}
			}
		}()

		authData, httpResp, ok := State.Authorize(r, req)

		if !ok {
			resp <- httpResp
			return
		}

		rd := &RouteData{
			Context: ctx,
			Auth:    authData,
		}

		if State.RouteDataMiddleware != nil {
			var err error
			rd, err = State.RouteDataMiddleware(rd, req)

			if err != nil {
				resp <- HttpResponse{
					Status: http.StatusInternalServerError,
					Json:   State.DefaultResponder.New(err.Error(), nil),
				}
				return
			}
		}

		resp <- r.Handler(*rd, req)
	}()

	respond(ctx, w, resp)
}

func marshalReq(r *http.Request, dst interface{}) (resp HttpResponse, ok bool) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		State.Logger.Error("[uapi/marshalReq] Failed to read body", zap.Error(err), zap.Int("size", len(bodyBytes)))
		return DefaultResponse(http.StatusInternalServerError), false
	}

	if len(bodyBytes) == 0 {
		return HttpResponse{
			Status: http.StatusBadRequest,
			Data:   State.Constants.BodyRequired,
		}, false
	}

	err = jsonimpl.Unmarshal(bodyBytes, &dst)

	if err != nil {
		State.Logger.Error("[uapi/marshalReq] Failed to unmarshal JSON", zap.Error(err), zap.Int("size", len(bodyBytes)))
		return HttpResponse{
			Status: http.StatusBadRequest,
			Json: State.DefaultResponder.New("Invalid JSON", map[string]string{
				"error": err.Error(),
			}),
		}, false
	}

	return HttpResponse{}, true
}

func MarshalReq(r *http.Request, dst any) (resp HttpResponse, ok bool) {
	return marshalReq(r, dst)
}

func MarshalReqWithHeaders(r *http.Request, dst any, headers map[string]string) (resp HttpResponse, ok bool) {
	resp, err := marshalReq(r, dst)

	resp.Headers = headers

	return resp, err
}
