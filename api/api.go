package api

import (
	"html/template"
	"net/http"

	"Chlamydia/config"
	"Chlamydia/constants"
	"Chlamydia/state"
	"Chlamydia/types"

	"Chlamydia/routes/primary"

	"github.com/go-chi/chi/v5"
	docs "github.com/infinitybotlist/eureka/doclib"
	"github.com/infinitybotlist/eureka/jsonimpl"
	"github.com/infinitybotlist/eureka/uapi"
	"go.uber.org/zap"

	_ "embed"
)

//go:embed docs.html
var docsHTML string
var openapi []byte

// Simple middleware to handle CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 50*1024*1024)

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "X-Client, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "X-Session-Invalid, Retry-After")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")

		if r.Method == "OPTIONS" {
			w.Write([]byte{})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func StartAPI() {
	// Create config
	config := config.NewConfig()

	// Setup Documentation
	docs.DocsSetupData = &docs.SetupData{
		URL:         "http://localhost:" + config.Port + "/",
		ErrorStruct: types.ApiError{},
		Info: docs.Info{
			Title:       config.Name,
			Version:     config.Version,
			Description: config.Description,
			Contact: docs.Contact{
				Name:  config.Contact.Name,
				URL:   config.Contact.URL,
				Email: config.Contact.Email,
			},
			License: docs.License{
				Name: "AGPL-3.0",
				URL:  config.GetLinkByName("License").URL,
			},
		},
	}
	docs.Setup()

	// Setup UAPI
	Setup()

	// Create router and apply middleware
	router := chi.NewRouter()
	router.Use(corsMiddleware)

	// Register API routes and services
	services := []uapi.APIRouter{
		primary.Router{},
	}

	for _, service := range services {
		name, desc := service.Tag()
		if name != "" {
			docs.AddTag(name, desc)
			uapi.State.SetCurrentTag(name)
		} else {
			panic("Woah! Service tags cannot be empty. Please fill it for usage.")
		}

		service.Routes(router)
	}

	// Serve documentation
	router.Get("/openapi", func(w http.ResponseWriter, r *http.Request) {
		w.Write(openapi)
	})

	// Handle errors correctly
	var err error

	// Load openapi to prevent large marshalling in all requests.
	openapi, err = jsonimpl.Marshal(docs.GetSchema())
	if err != nil {
		panic(err)
	}

	// Create docs template
	docsTempl := template.Must(template.New("docs").Parse(docsHTML))
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		docsTempl.Execute(w, map[string]string{
			"url": "/openapi",
		})
	})

	// Serve constant errors.
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(constants.EndpointNotFound))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(constants.MethodNotAllowed))
	})

	// Start Server
	state.Logger.Info("Chlamydia API started and accessible at http://localhost:" + config.Port + "/. Press CTRL+C to stop process.")
	cum := http.ListenAndServe(":"+config.Port, router)
	if cum != nil {
		state.Logger.Error("Failed to bind with socket: ", zap.Error(err))
	}
}
