package controllers

import (
	"fmt"
	"net/http"
)

//TestPage func
func TestPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
	<html lang="en">
	  <head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<meta http-equiv="X-UA-Compatible" content="ie=edge" />
		<title>Go WebSocket Tutorial</title>
	  </head>
	  <body>
		<h2>Hello World</h2>
	
		<script>
			let socket = new WebSocket("ws://localhost/arena");
			console.log("Attempting Connection...");
	
			socket.onopen = () => {
				console.log("Successfully Connected");
				socket.send("Hi From the Client!")
			};
			
			socket.onclose = event => {
				console.log("Socket Closed Connection: ", event);
				socket.send("Client Closed!")
			};
	
			socket.onerror = error => {
				console.log("Socket Error: ", error);
			};
	
		</script>
	  </body>
	</html>`

	w.Write([]byte(fmt.Sprintf(html)))
}
