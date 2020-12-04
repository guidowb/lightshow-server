package lightshow

import (
	"io"
	"log"
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
	router.HandleFunc("/connect", api.connect)
	return api
}

func (context api) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context.router.ServeHTTP(response, request)
}

func (context *api) getPattern(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	io.WriteString(response, context.pattern)
}

func (context *api) connect(response http.ResponseWriter, request *http.Request) {

	context.upgrader.CheckOrigin = func(r *http.Request) bool {
		if r.Header["Origin"][0] == "lightshow-arduino" {
			return true
		}
		log.Print("--- Refused access")
		log.Print("    Host: ", r.Host)
		log.Print("    URL: ", r.URL.String())
		headers := r.Header
		for key, value := range headers {
			log.Print("    ", key, ": ", value)
		}
		log.Print("---")
		return false
	}
	connection, err := context.upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Print("Upgrade error:", err)
		http.Error(response, "Upgrade Error", http.StatusBadRequest)
		return
	}
	defer connection.Close()

	log.Println("Client connected")

	done := make(chan struct{})
	go context.sendUpdates(connection, done)
	context.handleRequests(connection)
	close(done)
}

func (context *api) handleRequests(connection *websocket.Conn) {
	for {
		var message string
		err := connection.ReadJSON(&message)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Websocket error: %v", err)
			}
			return
		}
	}
}

func (context *api) sendUpdates(connection *websocket.Conn, done chan struct{}) {

	for { // Outer loop repeats when device pairing changes
		select {
		// case deviceUpdate := <-deviceChannel:
		// 	connection.WriteJSON(deviceUpdate)
		// 	break
		// case clientUpdate := <-clientChannel:
		// 	pairedDeviceID = clientUpdate.Get("pairedDeviceID").(string)
		// 	break
		case <-done:
			return
		}
	}
}
