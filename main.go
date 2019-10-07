package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	sw "./api"
	database "./api/database"
	ws "./api/websocket"
)

// Main function
func main() {
	hub := ws.NewHub()
	go hub.Run()
	flag.Parse()
	fmt.Println("Start Api..")
	//Init Database
	database.GetDatabase()
	//Init Router
	router := sw.NewRouter()
	//Start Webapi
	log.Fatal(http.ListenAndServe(":8080", router))

}
