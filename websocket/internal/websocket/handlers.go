package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type AuthRequest struct {
	Scheme      string `json:"scheme"`
	Credentials string `json:"credentials"`
}

type AuthResponse struct {
	Verify   bool   `json:"verify"`
	Id       int    `json:"id"`
	UserName string `json:"username"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func verifyTokenWithFastAPI(token string) (*AuthResponse, error) {
	fmt.Println("TOKEN IS:::::::::::::::::::::::")
	fmt.Println(token)
	authReq := AuthRequest{
		Scheme:      "bearer",
		Credentials: token,
	}
	jsonData, err := json.Marshal(authReq)
	fmt.Println("PASS jsonData, err := json.Marshal(authReq)")
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(
		"https://api.fantacytype.top/auth/verify",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	fmt.Println("PASS  calling api")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &result, nil
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	var username string
	protocols := websocket.Subprotocols(r)
	for _, protocol := range protocols {
		if strings.HasPrefix(protocol, "auth_token:") {
			token := strings.TrimPrefix(protocol, "auth_token:")
			auth, err := verifyTokenWithFastAPI(token)
			fmt.Println("GONNA SEE SOME ERRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR??")
			if err != nil {
				fmt.Println("here is the VALLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLIDDDDDDDDDD")
				fmt.Println(auth.Verify)
				username = auth.UserName
			}
		} else if strings.HasPrefix(protocol, "nickname:") {
			if username == "" {
				username = strings.TrimPrefix(protocol, "nickname:")
			}
		}
	}
	username = fmt.Sprintf("Guest_%d", time.Now().UnixNano()%10000)
	log.Printf("User comes as %s ", username)
	client := &Client{
		Conn:     conn,
		Id:       fmt.Sprintf("%s_%d", username, time.Now().UnixNano()),
		Hub:      hub,
		Username: username,
		Room:     r.URL.Query().Get("room"),
		SendChan: make(chan []byte, 256),
	}

	joinMessage := struct {
		Type     string `json:"type"`
		Username string `json:"username"`
		Content  string `json:"content"`
	}{
		Type:     "join",
		Username: username,
		Content:  fmt.Sprintf("%s joined the game", client.Username),
	}

	joinMessageBytes, _ := json.Marshal(joinMessage)
	client.broadcastToRoom(joinMessageBytes)
	//client.Hub.Register <- client
	fmt.Println("Client goes to client.hub.Register")

	go client.ReadPump()
	go client.WritePump()
}
