package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)






type Hub struct {
	Clients map[string]*Client
	Rooms map[string]map[string]*Client
	Games map[string]*GameState
	Register chan *Client
	Unregister chan *Client
	Mutex sync.Mutex
	Upgrader websocket.Upgrader
	ApiURL string
}



func NewHub(apiURL string) *Hub{
	rooms := make(map[string]map[string]*Client)
	rooms["English"] = make(map[string]*Client)
	rooms["English(2)"] = make(map[string]*Client)
	rooms["Persian"] = make(map[string]*Client)
	return &Hub{
		Clients: make(map[string]*Client),
		Rooms: rooms,
		Games: make(map[string]*GameState),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool	{
				return true
			},
			Subprotocols: []string{"auth_token", "nickname"},
		},
		ApiURL: apiURL,
	}
}


func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            // Lock here
			h.Mutex.Lock()
            h.Clients[client.Id] = client
            if _, exists := h.Rooms[client.Room]; !exists {
                h.Rooms[client.Room] = make(map[string]*Client)
            }
            h.Rooms[client.Room][client.Id] = client
			h.Mutex.Unlock()
			fmt.Println("HERE IS THE HUB HUB HUB CLIENTS")
			fmt.Println(h.Clients)
			// Unlock here
        case client := <-h.Unregister:
            h.Mutex.Lock()
            if _, ok := h.Clients[client.Id]; ok {
				if game , exists := h.Games[client.Room]; exists {
					delete(h.Rooms[client.Room], client.Id)
					delete(h.Clients, client.Id)
					delete(game.InGameUsers , client.Id)
				}else {
					delete(h.Rooms[client.Room], client.Id)
					delete(h.Clients, client.Id)
				}
                
                close(client.SendChan)
            }
            h.Mutex.Unlock()
        }
    }
}
