package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/parsaimi/elevenfinger_websocket/services"
)

type Client struct {
	Conn     *websocket.Conn `json:"-"`
	Hub      *Hub            `json:"-"`
	Id       string
	Username string
	Room     string
	SendChan chan []byte `json:"-"`
	IsReady  bool
}

type GameState struct {
	Text           string                       `json:"text"`
	StartTime      int64                        `json:"startTime"`
	EndTime        int64                        `json:"endGame"`
	IsActive       bool                         `json:"isActive"`
	PlayerProgress map[string]*PlayerWordRecord `json:"playerProgress"` // tracks words completed by each player
	leaderBoard    map[string]string         `json:"leaderBoard,omitempty"`
	TotalWords     int                          `json:"totalWords"`
	InGameUsers    map[string]*Client           `json:"inGameUsers"`
	wordList       []string                     `json:"wordList"`
	language       string                       `json:"language"`
	MatchId        string                       `json:"matchId"`
}

// readPump reads messages from websocket and pump it to the Hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
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
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.SendChan:
			fmt.Println("HERE IS THE USERNAME")
			fmt.Println(c.Username)
			fmt.Println("HERE IS THE ID")
			fmt.Println(c.Id)
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("error writing message:", err)
				return
			}
			log.Println(message)
			log.Println(string(message), "print in writePum we have the message")
		}
	}
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

	fmt.Println("rooms struct::::::::::")
	fmt.Println(roomsStatus)
	messageBytes, _ := json.Marshal(roomsStatus)
	fmt.Println("BYTEEEEEEEEEEEEEES in ROOOOOOOOMS status")
	fmt.Println(messageBytes)
	select {
	case c.SendChan <- messageBytes:
		fmt.Println("message goes to user in roomsStatus")
		fmt.Println(roomsStatus)
	default:
		close(c.SendChan)
		delete(c.Hub.Clients, c.Id)
	}

}

func (c *Client) roomStatus() {
	fmt.Println("ROMSSSSSSSSSSSSS STATUSSSSSSSSSSS ^^^^^^&&&&&&&&&&&&&****************")
	guests := make(map[string]bool)
	for _, value := range c.Hub.Rooms[c.Room] {
		guests[value.Username] = value.IsReady
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
		fmt.Println("IF LEVEL 1 PRINTING CLIENTS IN THAT ROOM")
		fmt.Println(clients)
		for _, thatclient := range clients {
			fmt.Println("IN LOOOOOOOOOOOOOOOOOOOOOOP")
			if _, gameExist := c.Hub.Games[c.Room]; gameExist {
				fmt.Println("IF LEVEL 2")
				if _, ok := c.Hub.Games[c.Room].InGameUsers[thatclient.Id]; !ok {
					fmt.Println("IF LEVEL 3")
					select {
					case thatclient.SendChan <- messageBytes:
						fmt.Println("message goes to user in first RommSTATUS :(((((((((((((((((((())))))))))))))))))))")
					default:
						close(c.SendChan)
						delete(c.Hub.Clients, c.Id)
					}

				}
			} else {
				fmt.Println("ELESE IS FINEEEEEEEEEEEEEEEEEE")
				select {
				case thatclient.SendChan <- messageBytes:
					fmt.Println("Message goes to user in RommStatus :))))))))))))")
				default:
					close(c.SendChan)
					delete(c.Hub.Clients, c.Id)
				}
			}

		}

		fmt.Println("AFTER LOOOOOOOOOOOOP")
	}

}

func (c *Client) joinPlayer(message json.RawMessage) {
	fmt.Println("JOIN TO ROOM HAPENS")
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
	c.Hub.Clients[c.Id] = c
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
	fmt.Println("HERE IS THE CLIENTTTTTTTTTSSSSSSSSSSSSSSS")
	fmt.Println(c.Hub.Clients)

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
		client.IsReady = false
	}
	inGameUsers := make(map[string]*Client)
	for id, client := range c.Hub.Rooms[c.Room] {
		inGameUsers[id] = client
	}

	language := "en"
	if c.Room == "Persian" {
		language = "fa"
	}

	displayText, wordList := services.GenerateCompetitionText(c.Room, language)
	gameState := &GameState{
		Text:           displayText,
		StartTime:      time.Now().UTC().Add(5 * time.Second).UnixMilli(),
		IsActive:       true,
		PlayerProgress: make(map[string]*PlayerWordRecord),
		leaderBoard:    make(map[string]string),
		TotalWords:     len(wordList),
		wordList:       wordList,
		language:       language,
		InGameUsers:    inGameUsers,
		MatchId:        uuid.New().String(),
	}
	for key, value := range c.Hub.Clients {
		gameState.PlayerProgress[key] = &PlayerWordRecord{
			Username:      value.Username,
			RemainedWords: wordList,
			StartTime:     time.Now().UTC().UnixMilli(),
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
	userWordInGame := &c.Hub.Games[c.Room].PlayerProgress[c.Id].RemainedWords
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
	for key, _ := range c.Hub.Rooms[c.Room] {
		c.Hub.Rooms[c.Room][key].IsReady = false
	}
	if len(c.Hub.Games[c.Room].leaderBoard) == len(c.Hub.Rooms[c.Room]) {

		endGameMessage := struct {
			Type        string             `json:"type"`
			LeaderBoard map[string]string `json:"leaderboard"`
		}{
			Type:        "endGame",
			LeaderBoard: c.Hub.Games[c.Room].leaderBoard,
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
					Type        string             `json:"type"`
					LeaderBoard map[string]string `json:"leaderboard"`
				}{
					Type: "endGame",
					LeaderBoard: c.Hub.Games[c.Room].leaderBoard,
				}
				messageBytes, _ := json.Marshal(endGameMessage)
				if c.Hub.Games[c.Room].MatchId == matchId {
					for _, playerInGame := range c.Hub.Games[c.Room].InGameUsers {
						select {
						case playerInGame.SendChan <- messageBytes:
						default:
							close(playerInGame.SendChan)
							delete(c.Hub.Clients, playerInGame.Id)
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

	endTime := time.Now().UTC().UnixMilli()
	startTime := c.Hub.Games[c.Room].PlayerProgress[c.Id].StartTime
	durationMs := endTime - startTime
	durationMinutes := float64(durationMs) / (1000.0 * 60.0)
	words := float64(len(c.Hub.Games[c.Room].Text)) / 5.0
	wpm := words / durationMinutes
	fmt.Println("HERE IS THE WWWWWWWWWPPPPPPPPMMMMMMMMMMM")
	fmt.Println(wpm)
	fmt.Println(durationMinutes)
	fmt.Println(len(c.Hub.Games[c.Room].Text))
	playerPosition := len(c.Hub.Games[c.Room].leaderBoard) + 1
	c.Hub.Games[c.Room].leaderBoard[strconv.Itoa(playerPosition)] = c.Username

	player := make(map[string]int)
	player[c.Id] = playerPosition
	playerRank := struct {
		Type            string  `json:"type"`
		PlayerId        string  `json:"playerid"`
		Rank            int     `json:"rank"`
		Wpm             float64 `json:"wpm"`
		DurationMinutes float64 `json:"durationminutes"`
	}{
		Type:            "playerRank",
		PlayerId:        c.Id,
		Rank:            playerPosition,
		Wpm:             wpm,
		DurationMinutes: durationMinutes,
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
	var auth struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}
	err := json.Unmarshal(messageContent, &auth)
	if err != nil {
		fmt.Println("error while parsing message in wordComplete method")
		return
	}
	if auth.Token != "" {
		fmt.Println("HERE IS TOKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKENKKN@@@@@@@@@@@@@@@@")
		fmt.Println(auth.Token)
		userData, err := verifyTokenWithFastAPI(auth.Token)
		if err != nil {
			fmt.Println("SOMTHING GOES WRONG WITH verifyTokenWithFastAPI")
		}
		fmt.Println("YOOOOOOOYOYOYOYOYOYOYOY")
		c.Username = userData.UserName
	} else {
		c.Username = auth.Username
	}

}

func (c *Client) joinRunningGame() {

	c.Hub.Games[c.Room].InGameUsers[c.Id] = c

	c.Hub.Games[c.Room].PlayerProgress[c.Id] = &PlayerWordRecord{
		Username:      c.Username,
		RemainedWords: c.Hub.Games[c.Room].wordList,
		StartTime:     time.Now().UTC().UnixMilli(),
	}

	startMessage := struct {
		Type      string   `json:"type"`
		Text      string   `json:"text"`
		Words     []string `json:"words"`
		StartTime int64    `json:"startTime"`
		Language  string   `json:"language"`
	}{
		Type:      "startGame",
		Text:      c.Hub.Games[c.Room].Text,
		Words:     c.Hub.Games[c.Room].wordList,
		StartTime: c.Hub.Games[c.Room].StartTime,
		Language:  c.Hub.Games[c.Room].language,
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
			case client.SendChan <- message:
			default:
				close(c.SendChan)
				delete(c.Hub.Rooms[c.Room], c.Id)
				delete(c.Hub.Clients, c.Id)
			}
		}
	}
}
