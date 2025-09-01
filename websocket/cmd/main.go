package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/parsaimi/elevenfinger_websocket/internal/websocket"
)


type AuthRequest struct {
	Scheme string `json:"scheme"`
	Credentials string `json:"credentials"`
}

type AuthResponse struct {
	Valid bool
	UserName string
}



func HealthCheck(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w, "Server is okay \n")
}


func main() {


	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	// create an instance of hub
    gameHub:= ws.NewHub("ws://127.0.0.1:8000/ws")

	// WAIT FOR REGISTER and UNREGISTER ( make client in hub , it makes users in rooms and clients )
    go gameHub.Run()

	// handling websocket connection starts here
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		ws.HandleWebSocket(gameHub , w , r)
	}) 
	log.Printf("Server starting on port %v", portString)
    log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}


