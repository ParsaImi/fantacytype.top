package hub

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/parsaimi/elevenfinger_websocket/internal/client"
)






type Hub struct {
	Clients map[string]*client.Client
	Rooms map[string]map[string]*client.Client
	Games map[string]*client.GameState
	Register chan *client.Client
	Unregister chan *client.Client
	Mutex sync.Mutex
	Upgrader websocket.Upgrader
	ApiURL string
}


func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            // Lock here
			h.Mutex.Lock()
            h.Clients[client.Id] = client
            if _, exists := h.Rooms[client.Room]; !exists {
                h.Rooms[client.Room] = make(map[string]*client.Client)
            }
            h.Rooms[client.Room][client.Id] = client
			h.Mutex.Unlock()
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
