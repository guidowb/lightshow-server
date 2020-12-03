package lightshow

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// API interface
type API interface {
	ServeHTTP(response http.ResponseWriter, request *http.Request)
}

type api struct {
	router   *mux.Router
	upgrader websocket.Upgrader
	pattern  string
}

// NewAPI returns an initialized Client API context
func NewAPI() API {
	router := mux.NewRouter().StrictSlash(true)
	api := api{
		router,
		websocket.Upgrader{},
		"solid #000",
	}
	router.HandleFunc("/pattern", api.getPattern)
	return api
}

func (context api) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context.router.ServeHTTP(response, request)
}

func (context *api) getPattern(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	io.WriteString(response, context.pattern)
}
