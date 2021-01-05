package main

import (
	lightshow "lightshow/api"
	"log"
	"net/http"
)

func main() {

	api := lightshow.NewAPI()

	log.Println("Listening")
	log.Fatal(http.ListenAndServe(":8080", api))
}
