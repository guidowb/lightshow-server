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

type clientinfo struct {
	MacAddress string
	Program    string
	Connection *websocket.Conn
	Commands   chan string
}

type api struct {
	router   *mux.Router
	upgrader websocket.Upgrader
	clients  map[string]*clientinfo
	programs map[string]string
}

// NewAPI returns an initialized Client API context
func NewAPI() API {
	router := mux.NewRouter().StrictSlash(true)
	api := api{
		router,
		websocket.Upgrader{},
		map[string]*clientinfo{},
		map[string]string{},
	}
	router.HandleFunc("/pattern/{program}", api.setPattern).Methods(http.MethodPost)
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

	vars := mux.Vars(request)
	program := vars["program"]
	log.Printf("For program: %s\n", program)

	context.programs[program] = pattern

	for _, client := range context.clients {
		if client.Program == program {
			log.Printf("Updating client: %s\n", client.MacAddress)
			client.Commands <- "pattern\n" + pattern
		}
	}
	response.WriteHeader(http.StatusCreated)
}

func (context *api) getCookie(request *http.Request, name string) *string {

	cookies := request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return &cookie.Value
		}
	}
	return nil
}

func (context *api) connect(response http.ResponseWriter, request *http.Request) {

	macAddress := context.getCookie(request, "macAddress")
	program := context.getCookie(request, "program")

	context.upgrader.CheckOrigin = func(r *http.Request) bool {
		if macAddress != nil && program != nil && r.Header["Origin"][0] == "lightshow-arduino" {
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

	client, present := context.clients[*macAddress]
	if present {
		log.Printf("Client %s reconnected, program %s\n", *macAddress, *program)
		close(client.Commands)
		client.Connection.Close()
	} else {
		log.Printf("Client %s connected, program %s\n", *macAddress, *program)
	}
	client = &clientinfo{
		*macAddress,
		*program,
		connection,
		make(chan string),
	}
	context.clients[*macAddress] = client

	go context.sendCommands(client)
	if pattern, ok := context.programs[*program]; ok {
		client.Commands <- "pattern\n" + pattern
	}
	context.handleUpdates(client)

	delete(context.clients, *macAddress)
	log.Printf("Client %s disconnected\n", *macAddress)
}

func (context *api) handleUpdates(client *clientinfo) {
	for {
		_, message, err := client.Connection.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Client %s websocket error: %v", client.MacAddress, err)
			}
			return
		}
		log.Printf("Client %s sent message: %s\n", client.MacAddress, message)
	}
}

func (context *api) sendCommands(client *clientinfo) {
	for {
		select {

		case <-time.After(35 * time.Second):
			// log.Printf("Client %s getting a ping\n", client.MacAddress)
			client.Connection.WriteControl(websocket.PingMessage, []byte{}, time.Time{})

		case command, ok := <-client.Commands:
			if ok {
				log.Printf("Client %s getting command: %s\n", client.MacAddress, command)
				client.Connection.WriteMessage(websocket.TextMessage, []byte(command))
			} else {
				return
			}
		}
	}
}
