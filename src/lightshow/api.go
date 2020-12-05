package lightshow

import (
	"bytes"
	"log"
	"net/http"
	"time"

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
	commands chan string
}

// NewAPI returns an initialized Client API context
func NewAPI() API {
	router := mux.NewRouter().StrictSlash(true)
	api := api{
		router,
		websocket.Upgrader{},
		"solid #000",
		make(chan string),
	}
	router.HandleFunc("/pattern", api.setPattern).Methods(http.MethodPost)
	router.HandleFunc("/connect", api.connect)
	return api
}

func (context api) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context.router.ServeHTTP(response, request)
}

func (context *api) setPattern(response http.ResponseWriter, request *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(request.Body)
	pattern := buf.String()
	log.Println("Received pattern:")
	log.Println(pattern)
	context.commands <- "pattern\n" + pattern
	response.WriteHeader(http.StatusCreated)
}

func (context *api) getCookie(request *http.Request, name string) *http.Cookie {

	cookies := request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
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

	macAddress := context.getCookie(request, "macAddress")
	connection, err := context.upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Print("Upgrade error:", err)
		http.Error(response, "Upgrade Error", http.StatusBadRequest)
		return
	}
	defer connection.Close()

	log.Println("Client connected, ", macAddress)
	connection.WriteControl(websocket.PingMessage, []byte{}, time.Time{})

	done := make(chan struct{})
	go context.sendUpdates(connection, done)
	context.handleRequests(connection)
	close(done)

	log.Println("Client disconnected, ", macAddress)
}

func (context *api) handleRequests(connection *websocket.Conn) {
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Websocket error: %v", err)
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Websocket error: %v", err)
			}
			return
		}
		log.Printf("Received message")
		log.Printf("Type: %s", messageType)
		log.Printf("Data: %s", message)
	}
}

func (context *api) sendUpdates(connection *websocket.Conn, done chan struct{}) {

	for { // Outer loop repeats when device pairing changes
		select {
		case <-time.After(35 * time.Second):
			log.Printf("Sending a ping")
			connection.WriteControl(websocket.PingMessage, []byte{}, time.Time{})
			break
		case command := <-context.commands:
			log.Printf("Sending command: %s\n", command)
			connection.WriteMessage(websocket.TextMessage, []byte(command))
			break
		case <-done:
			return
		}
	}
}
