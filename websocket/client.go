package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/parsaimi/elevenfinger_websocket/internal/hub"
	"github.com/parsaimi/elevenfinger_websocket/internal/message"
	"github.com/parsaimi/elevenfinger_websocket/services"
)

type Client struct {
	Conn     *websocket.Conn
	Hub      *hub.Hub
	Id       string
	Username string
	Room     string
	SendChan chan []byte
	IsReady  bool
}

type GameState struct {
	Text           string                       `json:"text"`
	StartTime      int64                        `json:"startTime"`
	IsActive       bool                         `json:"isActive"`
	PlayerProgress map[string]*message.PlayerWordRecord `json:"playerProgress"` // tracks words completed by each player
	leaderBoard    map[string]*Client           `json:"leaderBoard,omitempty"`
	TotalWords     int                          `json:"totalWords"`
	InGameUsers    map[string]*Client           `json:"inGameUsers"`
	wordList       []string                     `json:"wordList"`
	language       string                       `json:"language"`
	MatchId        string                       `json:"matchId"`
}

// readPump reads messages from websocket and pump it to the Hub
func (c *Client) ReadPump() {

	for {
		_, message, error := c.Conn.ReadMessage()
		if error != nil {
			log.Printf("error : %v", error)
			break
		}

		c.handleGameMessage(message)

	}

}

func (c *Client) WritePump() {

}

func (c *Client) handleGameMessage(message []byte) {
	var gameMessage struct {
		Type    string          `json:"type"`
		Content json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(message, &gameMessage); err != nil {
		log.Printf("Error parsing message")
		return
	}
	switch gameMessage.Type {
	case "roomsStatus":
		c.roomsStatus()
	case "join":
		c.joinPlayer(gameMessage.Content)
	case "ready":
		c.readyPlayer()
	case "startGame":
		c.readyToStart()
	case "wordComplete":
		c.wordComplete(gameMessage.Content)
	case "endGame":
		log.Println("+++++++++++++ENG GAME SOLDIER+++++++++++++")
	case "roomStatus":
		log.Println("going into roomStatus Method ( in handleGameMessage)")
		c.roomStatus()
	case "usercred":
		c.userCred(gameMessage.Content)

	}

}

func (c *Client) roomsStatus() {
	roomsStatus := struct {
		Type  string                        `json:"type"`
		Rooms map[string]map[string]*Client `json:"rooms"`
	}{
		Type:  "roomsStatus",
		Rooms: c.Hub.Rooms,
	}

	messageBytes, _ := json.Marshal(roomsStatus)
	select {
	case c.SendChan <- messageBytes:
		fmt.Println("message goes to user")
	default:
		close(c.SendChan)
		delete(c.Hub.Clients, c.Id)
	}

}

func (c *Client) roomStatus() {

	guests := make(map[string]bool)
	for _, value := range c.Hub.Rooms[c.Room] {
		guests[value.username] = value.isReady
	}
	roomStatus := struct {
		Type    string          `json:"type"`
		Players map[string]bool `json:"players"`
	}{
		Type:    "roomStatus",
		Players: guests,
	}

	messageBytes, _ := json.Marshal(roomStatus)
	if clients, ok := c.Hub.Rooms[c.Room]; ok {
		for _, thatclient := range clients {
			if _, gameExist := c.Hub.Games[c.Room]; gameExist {
				if _, ok := c.Hub.Games[c.Room].InGameUsers[thatclient.id]; !ok {
					select {
					case thatclient.sendChan <- messageBytes:
					default:
						close(c.SendChan)
						delete(c.Hub.Clients, c.Id)
					}

				}
			} else {
				select {
				case thatclient.sendChan <- messageBytes:
				default:
					close(c.SendChan)
					delete(c.Hub.Clients, c.Id)
				}
			}

		}
	}

}

func (c *Client) joinPlayer(message json.RawMessage) {
	if len(c.Hub.Rooms[c.Room]) == 0 {
		go func() {
			<-time.After(10 * time.Second)
			c.readyToStart()
		}()
	}
	var result map[string]string
	err := json.Unmarshal(message, &result)
	if err != nil {
		fmt.Println("Error while parsing message in join player")
	}
	room := result["room"]
	c.Room = room
	c.Hub.Mutex.Lock()
	// c.hub.Clients[c.id] = c
	if _, exists := c.Hub.Rooms[c.Room]; !exists {
		c.Hub.Rooms[c.Room] = make(map[string]*Client)
	}
	c.Hub.Rooms[c.Room][c.Id] = c
	c.Hub.Mutex.Unlock()
	_, exist := c.Hub.Games[c.Room]
	if exist {
		if c.Hub.Games[c.Room].IsActive {
			c.joinRunningGame()
		} else {
			c.roomStatus()
		}
	} else {
		c.roomStatus()
	}

}

func (c *Client) readyPlayer() {
	c.IsReady = true
	c.roomStatus()
	c.readyToStart()
}

func (c *Client) readyToStart() {
	if _, exists := c.Hub.Games[c.Room]; !exists {
		c.startNewGame()
	} else if !c.Hub.Games[c.Room].IsActive {
		c.startNewGame()
	}

}

func (c *Client) startNewGame() {
	for _, client := range c.Hub.Rooms[c.Room] {
		client.isReady = false
	}
	displayText, wordList := services.GenerateCompetitionText(c.Room)
	inGameUsers := make(map[string]*Client)
	for id, client := range c.Hub.Rooms[c.Room] {
		inGameUsers[id] = client
	}
	language := "en"
	if c.Room == "room3" {
		language = "fa"
	}
	gameState := &GameState{
		Text:           displayText,
		StartTime:      time.Now().UTC().Add(5 * time.Second).UnixMilli(),
		IsActive:       true,
		PlayerProgress: make(map[string]*message.PlayerWordRecord),
		leaderBoard:    make(map[string]*Client),
		TotalWords:     len(wordList),
		wordList:       wordList,
		language:       language,
		InGameUsers:    inGameUsers,
		MatchId:        uuid.New().String(),
	}
	for key, value := range c.Hub.Clients {
		gameState.PlayerProgress[key] = &message.PlayerWordRecord{
			Username:      value.username,
			RemainedWords: wordList,
		}
	}
	c.Hub.Games[c.Room] = gameState

	startMessage := struct {
		Type     string   `json:"type"`
		Text     string   `json:"text"`
		Words    []string `json:"words"`
		Time     int64    `json:"startTime"`
		Language string   `json:"language"`
	}{
		Type:     "startGame",
		Text:     gameState.Text,
		Words:    wordList,
		Time:     gameState.StartTime,
		Language: language,
	}

	messageBytes, _ := json.Marshal(startMessage)
	c.broadcastToRoom(messageBytes)

}

func (c *Client) wordComplete(messageContent json.RawMessage) {
	var result map[string]string
	err := json.Unmarshal(messageContent, &result)
	if err != nil {
		fmt.Println("error while parsing message in wordComplete method")
	}
	userInputWord := result["word"]
	userWordInGame := &c.Hub.Games[c.Room].PlayerProgress[c.Id].remainedWords
	if userInputWord != (*userWordInGame)[0] {
		fmt.Println(userWordInGame)
		fmt.Println(userInputWord)
		fmt.Println("CHEATER SPOTTED !!!")
		return
	}
	totalWords := c.Hub.Games[c.Room].TotalWords
	*userWordInGame = (*userWordInGame)[1:]
	completedWords := totalWords - len(*userWordInGame)
	var userProgressBar int
	userProgressBar = int(math.Round((float64(completedWords) / float64(totalWords)) * 100))
	c.userProgress(userProgressBar)

	if len(*userWordInGame) == 0 {
		c.userRanking()
		c.endGame(c.Hub.Games[c.Room].MatchId)
	}

}

func (c *Client) endGame(matchId string) {
	for key, value := range c.Hub.Rooms[c.Room] {
		c.Hub.Rooms[c.Room][key].isReady = false
	}
	if len(c.Hub.Games[c.Room].leaderBoard) == len(c.Hub.Rooms[c.Room]) {

		endGameMessage := struct {
			Type string `json:"type"`
		}{
			Type: "endGame",
		}
		c.Hub.Games[c.Room].IsActive = false
		messageBytes, _ := json.Marshal(endGameMessage)
		c.broadcastToRoom(messageBytes)
		return

	}
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		timeRemaining := 20

		for range ticker.C {
			timeRemaining--

			if timeRemaining <= 0 {
				c.Hub.Games[c.Room].IsActive = false
				ticker.Stop()

				//gs.mutex.Lock()
				//gs.inProgress = false
				//gs.mutex.Unlock()
				
				endGameMessage := struct {
					Type string `json:"type"`
				}{
					Type: "endGame",
				}
				messageBytes, _ := json.Marshal(endGameMessage)
				if c.Hub.Games[c.Room].MatchId == matchId {
					for _, playerInGame := range c.Hub.Games[c.Room].InGameUsers {
						select {
						case playerInGame.sendChan <- messageBytes:
						default:
							close(playerInGame.sendChan)
							delete(c.Hub.Clients, playerInGame.id)
						}

					}
				}

				break
			}

			// Optional: broadcast time updates (every second or at intervals)
		}
	}()

}

func (c *Client) userRanking() {
	playerPosition := len(c.Hub.Games[c.Room].leaderBoard) + 1
	c.Hub.Games[c.Room].leaderBoard[strconv.Itoa(playerPosition)] = c

	player := make(map[string]int)
	player[c.Id] = playerPosition
	playerRank := struct {
		Type       string         `json:"type"`
		PlayerRank map[string]int `json:"playerrank"`
	}{
		Type:       "playerRank",
		PlayerRank: player,
	}

	messageBytes, _ := json.Marshal(playerRank)
	c.broadcastToRoom(messageBytes)

}

func (c *Client) userProgress(progress int) {
	progressMessage := struct {
		Type       string `json:"type"`
		Userid     string `json:"userid"`
		Percentage int    `json:"percentage"`
	}{
		Type:       "userProgress",
		Userid:     c.Username,
		Percentage: progress,
	}
	messageBytes, _ := json.Marshal(progressMessage)
	c.broadcastToRoom(messageBytes)

}

func (c *Client) userCred(messageContent json.RawMessage) {
	var result map[string]string
	err := json.Unmarshal(messageContent, &result)
	if err != nil {
		fmt.Println("error while parsing message in wordComplete method")
	}
	c.Username = result["username"]
}

func (c *Client) joinRunningGame() {

	c.Hub.Games[c.Room].InGameUsers[c.Id] = c

	c.Hub.Games[c.Room].PlayerProgress[c.Id] = &message.PlayerWordRecord{
		Username:      c.Username,
		RemainedWords: c.Hub.Games[c.Room].wordList,
	}

	startMessage := struct {
		Type     string   `json:"type"`
		Text     string   `json:"text"`
		Words    []string `json:"words"`
		Time     int64    `json:"startTime"`
		Language string   `json:"language"`
	}{
		Type:     "startGame",
		Text:     c.Hub.Games[c.Room].Text,
		Words:    c.Hub.Games[c.Room].wordList,
		Time:     c.Hub.Games[c.Room].StartTime,
		Language: c.Hub.Games[c.Room].language,
	}

	messageBytes, _ := json.Marshal(startMessage)

	select {
	case c.SendChan <- messageBytes:
		fmt.Println("message goes to user in the gameee")
	default:
		close(c.SendChan)
		delete(c.Hub.Clients, c.Id)
	}
}

func (c *Client) broadcastToRoom(message []byte) {
	c.Hub.Mutex.Lock()
	defer c.Hub.Mutex.Unlock()
	if clients, ok := c.Hub.Rooms[c.Room]; ok {
		for _, client := range clients {
			select {
			case client.sendChan <- message:
			default:
				close(c.SendChan)
				delete(c.Hub.Rooms[c.Room], c.Id)
				delete(c.Hub.Clients, c.Id)
			}
		}
	}
}
