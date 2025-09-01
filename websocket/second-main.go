package main

import (
	//"github.com/parsaimi/elevenfinger_websocket/internal/app"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/exp/rand"
)


type Client struct {
	conn *websocket.Conn
	id string
	username string
	room string
	sendChan chan []byte
	isReady bool
}


type SentenceResponse struct {
	Sentence string `json:"sentence"`
	Locale   string `json:"locale"`
}
type GameServer struct {
	clients map[string]*Client
	rooms map[string]map[string]*Client
	games map[string]*GameState
	register chan *Client
	unregister chan *Client
	mutex sync.Mutex
	upgrader websocket.Upgrader
	apiURL string
}

type PlayerWordRecord struct {
	username        string           
	remainedWords   []string	
}

type GameState struct {
	Text            string            `json:"text"`
	StartTime       int64         `json:"startTime"`
	IsActive        bool              `json:"isActive"`
	PlayerProgress  map[string]*PlayerWordRecord    `json:"playerProgress"` // tracks words completed by each player
	leaderBoard     map[string]*Client          `json:"leaderBoard,omitempty"`
	TotalWords      int              `json:"totalWords"`
	InGameUsers      map[string]*Client        `json:"inGameUsers"`
	wordList       []string         `json:"wordList"`
	language        string          `json:"language"`
	MatchId         string           `json:"matchId"`
}

type GameMessage struct {
	Type            string           `json:"type"`
}

type roomStatus struct {
	Type            string           `json:"type"`
	Players         map[string]bool
}


func NewGameServer(apiURL string) *GameServer {
	rooms := make(map[string]map[string]*Client)
	rooms["room1"] = make(map[string]*Client)
	rooms["room2"] = make(map[string]*Client)
	rooms["room3"] = make(map[string]*Client)
	return &GameServer{
		clients: make(map[string]*Client),
		rooms: rooms,
		games: make(map[string]*GameState),
		register: make(chan *Client),
		unregister: make(chan *Client),
		apiURL: apiURL,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool	{
				return true
			},
			Subprotocols: []string{"auth_token", "nickname"},
		},
	}
}

func (gs *GameServer) userRanking(client *Client){
	playerPosition := len(gs.games[client.room].leaderBoard) + 1
	gs.games[client.room].leaderBoard[strconv.Itoa(playerPosition)] = client

	player := make(map[string]int)
	player[client.id] = playerPosition
	playerRank := struct {
		Type     string           `json:"type"`
		PlayerRank  map[string]int  `json:"playerrank"`
	}{
		Type: "playerRank",
		PlayerRank: player,
	}
			
	messageBytes, _ := json.Marshal(playerRank)
	gs.broadcastToRoom(client.room , messageBytes)

}



func generateCompetitionText(room string) (string, []string ) {
	
	
	// Common Persian words for room3
	persianWords := []string{
		"سلام", "جهان", "کتاب", "خانه", "درخت", "آزادی", "عشق", "دوست", "خورشید", "ماه",
		"ستاره", "آسمان", "زمین", "رود", "کوه", "گل", "پرنده", "باغ", "شهر", "روستا",
	}

	// English words for other rooms
	englishWords := []string{
		"introductory", "preliminary", "ceremonial", "demonstrating", "graves", "speed", "hoses", "steep", "reunions", "I", "You", "Yeah",
		"script", "devoted", "prepositions", "indie", "fascinating", "courage", "Star", "Five", "outside",
	}
	
	baseurl := "http://127.0.0.1:8000/corpus/sentence"

	u , _ := url.Parse(baseurl)

	q := u.Query()
	q.Set("locale", "en_US")
	u.RawQuery = q.Encode()
	resp, _ := http.Get(u.String())
	defer resp.Body.Close()

	body, _:= io.ReadAll(resp.Body)
	fmt.Println(string(body), "CHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH")

	var response SentenceResponse
	err := json.Unmarshal(body , &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	splitedWords := strings.Split(response.Sentence, " ")

	fmt.Println(response.Sentence, "CHECKCHECKCHECKCHECKCHECK")

	var words []string
	var result []string

	if room == "room3" {
		words = persianWords
	} else {
		words = englishWords
	}

	// Check if the word list is empty
	if len(words) == 0 {
		fmt.Println("ohhhahahahahha")
	}

	// Generate 10 random words
	for i := 0; i < 10; i++ {
		word := words[rand.Intn(len(words))]
		result = append(result, word)
	}


	return response.Sentence, splitedWords
}

func (gs *GameServer) Run() {
    for {
        select {
        case client := <-gs.register:
            // Lock here
			gs.mutex.Lock()
            gs.clients[client.id] = client
            if _, exists := gs.rooms[client.room]; !exists {
                gs.rooms[client.room] = make(map[string]*Client)
            }
            gs.rooms[client.room][client.id] = client
			gs.mutex.Unlock()
			// Unlock here
        case client := <-gs.unregister:
            gs.mutex.Lock()
            if _, ok := gs.clients[client.id]; ok {
				if game , exists := gs.games[client.room]; exists {
					delete(gs.rooms[client.room], client.id)
					delete(gs.clients, client.id)
					delete(game.InGameUsers , client.id)
				}else {
					delete(gs.rooms[client.room], client.id)
					delete(gs.clients, client.id)
				}
                
                close(client.sendChan)
            }
            gs.mutex.Unlock()
        }
    }
}


func (gs *GameServer) roomsStatus(client *Client){
	roomsStatus := struct {
		Type    string                         `json:"type"` 
		Rooms   map[string]map[string]*Client  `json:"rooms"` 
  
	}{
		Type: "roomsStatus",
		Rooms: gs.rooms,
	}

	messageBytes, _ := json.Marshal(roomsStatus)
	select{
			case client.sendChan <- messageBytes:
				fmt.Println("message goes to user")	
			default:
				close(client.sendChan)
				delete(gs.clients, client.id)
			}


}

func (gs *GameServer) roomStatus(client *Client){
	
	fmt.Println(client.username, "username hereeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee!!!@!@!@!@!@!@!@!@!@!@@!@!")
	guests := make(map[string]bool)
	for key , value := range gs.rooms[client.room]{
		fmt.Println(key, "salam in key hasttttttttttttttttttttttttttttttttttttttttttttttttttttt")
		guests[value.username] = value.isReady
	}
	roomStatus := struct {
		Type     string           `json:"type"`
		Players  map[string]bool  `json:"players"`
	}{
		Type: "roomStatus",
		Players: guests,
	}
			
	messageBytes, _ := json.Marshal(roomStatus)
	if clients, ok := gs.rooms[client.room]; ok{
		for _, thatclient := range clients {
			if _ , gameExist := gs.games[client.room]; gameExist {
				fmt.Println(gs.games[client.room].InGameUsers ,  "***********************************************************************************************************************")
				if _ , ok := gs.games[client.room].InGameUsers[thatclient.id]; !ok {
					fmt.Println("where we want)()()()(((((((((((((((((((((((((((((((((()))))))))))))))))))))))))))))))))")
					select{
						case thatclient.sendChan <- messageBytes:
							fmt.Println("message goes to user")	
						default:
							close(client.sendChan)
							delete(gs.clients, client.id)
						}

				}
			}else{
					select{
						case thatclient.sendChan <- messageBytes:
							fmt.Println("message goes to user")	
						default:
							close(client.sendChan)
							delete(gs.clients, client.id)
						}
					}	


		}
	}


}

func (gs *GameServer) joinPlayer(client *Client, message json.RawMessage){
	if len(gs.rooms[client.room]) == 0 {
		go func() {
			<-time.After(10 * time.Second)
			gs.readyToStart(client.room)
		}()
	}
	fmt.Println(client.username, "eyesssssssssssssssssssssssssssssss on thissssssssssssssssssssssssssssssssssssss")
	var result map[string]string
	err := json.Unmarshal(message , &result)
	if err != nil {
		fmt.Println("Error while parsing message in join player")
	}
	room := result["room"]
	nickname := result["nickname"]
	fmt.Println("THIS IS NICKNAME JAKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK#################### ", nickname)
	fmt.Println(room , "here is the messeage !@")
	client.room = room
	gs.mutex.Lock()
	gs.clients[client.id] = client
    if _, exists := gs.rooms[client.room]; !exists {
        gs.rooms[client.room] = make(map[string]*Client)
    }
    gs.rooms[client.room][client.id] = client
	gs.mutex.Unlock()
	_ , exist := gs.games[client.room]
	if exist{
		fmt.Println("ROOOOOOOOOOOMMMMMMMMMMMMMMMMMMMMMMMMMMMmm    EXISSSSSSSSSSSSSSSSSSSTTTTTTTTTT")
		if gs.games[client.room].IsActive {
			fmt.Println("ROOOOOOOOOOOOOOOOOOOOOOOOOOOOM IS ACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCc!!!!!!!!!!!!!!!")
			gs.joinRunningGame(client)
			fmt.Println(gs.rooms[client.room], "IIIIIIIIIIIIIIIIIIIIIIIIIIIISSSSSSSSSSSSSSSS PLAYERSSSSSSSSSSSSSSSSSS(*#&#(&@&#(&#@@(()))))")
		}else{
			fmt.Println("ROOOOOOOOOOOOOOOOOOOOOOOOOOOOM IS NNNNNNNNOOOOOOOOOOOOOOTTTTTTTTTTTTT  ACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCc!!!!!!!!!!!!!!!")
			gs.roomStatus(client)
		}
	}else{
		fmt.Println("ROOOOOOOOOOOOOOOOOOOOOOOOOOOOM IS DOWSNT EXITS AT ALLL") 
		gs.roomStatus(client)
	}
	
	fmt.Println("here is the player room" , client.room)
	fmt.Println("important one ########################################" , gs.rooms[client.room])
	
	

}



func (gs *GameServer) readyPlayer(client *Client){
	client.isReady = true
	gs.roomStatus(client)
	for key, value := range gs.rooms[client.room] {
		fmt.Println(key, "")
        if value.isReady != true{
			fmt.Println("player ", value.username, "with this id : ", key ," is not ready and game does not going to start")
			return
		}else{
			fmt.Println(key , value.isReady , "yoooooooooo pretty motherfoucker gooooooooooooooooooooooooooooooo")
		}
		
    } 
	gs.readyToStart(client.room)
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAfter rommstatus", gs.games[client.room].InGameUsers)
}


func (gs *GameServer) endGame(client *Client, matchId string){
	fmt.Println("end offfffffffffffffffffffffff gameeeeeeeeeeeeeeeeeeeee")
	for key, value := range gs.rooms[client.room] {
		fmt.Println(value.isReady , "active brossss")
		gs.rooms[client.room][key].isReady = false
	}
	if len(gs.games[client.room].leaderBoard) == len(gs.rooms[client.room]){
			
			endGameMessage := struct{
				Type			  string				`json:"type"`
			}{
				Type: "endGame",

			}
		gs.games[client.room].IsActive = false
		fmt.Println("game activity goes ***&*&&&&&&&&&&&&&&*&*&*&*&*&&**&*%*%*%*%*%*%#*^#(@%&#%*&%#&*%&#%&^#*#&%@^&#$*@^*(#!!!!!!!!!!!!!!!!!!!!!!!!!!!!!))", gs.games[client.room].IsActive )
		messageBytes, _ := json.Marshal(endGameMessage)
		gs.broadcastToRoom(client.room , messageBytes)
		return


	}
	go func() {
        ticker := time.NewTicker(1 * time.Second)
        timeRemaining := 20
        
        for range ticker.C {
            timeRemaining--
            
            if timeRemaining <= 0 {
				gs.games[client.room].IsActive = false
                ticker.Stop()
                
                //gs.mutex.Lock()
                //gs.inProgress = false
                //gs.mutex.Unlock()
				endGameMessage := &GameMessage{
					Type: "endGame",
				}
				messageBytes, _ := json.Marshal(endGameMessage)
				if gs.games[client.room].MatchId == matchId{
					fmt.Println("room status is " , gs.games[client.room].IsActive)
					fmt.Println("players are ", gs.games[client.room].InGameUsers)
					for _ , playerInGame := range gs.games[client.room].InGameUsers {
						select{
							case playerInGame.sendChan <- messageBytes:
								fmt.Println("message goes to user in the gameee")	
							default:
								close(playerInGame.sendChan)
								delete(gs.clients, playerInGame.id)
							}
						
						
					}
			}

                
                break
            }
            
            // Optional: broadcast time updates (every second or at intervals)
        }
    }()	

}

func (gs *GameServer) userProgress(client *Client, progress int){
	progressMessage := struct {
		Type string      `json:"type"`
		Userid string  `json:"userid"`
		Percentage int   `json:"percentage"`

	}{
		Type: "userProgress",
		Userid: client.username,
		Percentage: progress,
	}
	messageBytes, _ := json.Marshal(progressMessage)
	gs.broadcastToRoom(client.room , messageBytes)



}

func (gs *GameServer) wordComplete(client *Client , messageContent json.RawMessage){
	fmt.Println("a word Complete message comes in ! ( print in wordComplete method )")
	var result map[string]string
	err := json.Unmarshal(messageContent, &result)
	if err != nil{
		fmt.Println("error while parsing message in wordComplete method")
	}
	fmt.Println(result)
	fmt.Println("WE ARE COLL99999999999999999999999999999999999")
	userInputWord := result["word"]
	fmt.Println("are we coll?88888888888888888888888888888888888")
	fmt.Println(gs.games[client.room].PlayerProgress, "{}{{{{{{{{{PLAYER PROGRESSJk}}}}}}}}}")
	userWordInGame := &gs.games[client.room].PlayerProgress[client.id].remainedWords
	if userInputWord != (*userWordInGame)[0]{
		fmt.Println(userWordInGame)
		fmt.Println(userInputWord)
		fmt.Println("CHEATER SPOTTED !!!")
		return
	}
	totalWords := gs.games[client.room].TotalWords
	*userWordInGame = (*userWordInGame)[1:]
	fmt.Println("i DID A SUBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBbb ", userWordInGame)
	completedWords := totalWords - len(*userWordInGame)
	var userProgressBar int
	userProgressBar = int(math.Round((float64(completedWords) / float64(totalWords)) * 100))
	fmt.Println("this the % of player progress", userProgressBar)
	fmt.Println(userWordInGame, ")()()()()()()()______________________")
	gs.userProgress(client , userProgressBar)
	
	if len(*userWordInGame) == 0{
		fmt.Println("user finished the game @@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		gs.userRanking(client)
		fmt.Println(gs.games[client.room].leaderBoard)
		gs.endGame(client , gs.games[client.room].MatchId)
	}
		
	

}


func (gs *GameServer) joinRunningGame(client *Client){

		gs.games[client.room].InGameUsers[client.id] = client

		gs.games[client.room].PlayerProgress[client.id] = &PlayerWordRecord{
			username: client.username,
			remainedWords: gs.games[client.room].wordList,
		}
		
		startMessage := struct {
			Type string    `json:"type"`
			Text string    `json:"text"`
			Words []string `json:"words"`
			Time int64 `json:"startTime"`
			Language string `json:"language"`
		}{
			Type: "startGame",
			Text: gs.games[client.room].Text,
			Words: gs.games[client.room].wordList,
			Time: gs.games[client.room].StartTime,
			Language: gs.games[client.room].language,
		}	

		messageBytes, _ := json.Marshal(startMessage)

		select{
			case client.sendChan <- messageBytes:
				fmt.Println("message goes to user in the gameee")	
			default:
				close(client.sendChan)
				delete(gs.clients, client.id)
		}
}


func (gs *GameServer) readyToStart(room string){
	if _, exists := gs.games[room]; !exists {
		gs.startNewGame(room)
	}else if !gs.games[room].IsActive {
		gs.startNewGame(room)
	}

}

func (gs *GameServer) startNewGame(room string){
		for _, client := range gs.rooms[room]{
			client.isReady = false
		}
		fmt.Println("game is going to start in this room!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!" , gs.rooms[room])	
		displayText , wordList := generateCompetitionText(room)
		inGameUsers := make(map[string]*Client)
		for id, client := range gs.rooms[room] {
			inGameUsers[id] = client
		}
		language := "en"
		if room == "room3" {
			language = "fa"
		}
		gameState := &GameState{
			Text:           displayText,
			StartTime:      time.Now().UTC().Add(5 * time.Second).UnixMilli(),
			IsActive:       true,
			PlayerProgress: make(map[string]*PlayerWordRecord),
			leaderBoard:    make(map[string]*Client),
			TotalWords:     len(wordList),
			wordList:       wordList,
			language:       language,
			InGameUsers:    inGameUsers,
			MatchId:		uuid.New().String(),
		}
		fmt.Println("the game should start at ", gameState.StartTime , "DRE")
		for key , value := range gs.clients {
			gameState.PlayerProgress[key] = &PlayerWordRecord{
				username: value.username,
				remainedWords: wordList,
			}
		}
		gs.games[room] = gameState

		
		

		startMessage := struct {
			Type string    `json:"type"`
			Text string    `json:"text"`
			Words []string `json:"words"`
			Time int64 `json:"startTime"`
			Language string `json:"language"`
		}{
			Type: "startGame",
			Text: gameState.Text,
			Words: wordList,
			Time: gameState.StartTime,
			Language: language,
		}
		log.Println("in start New Game Line 129")
		for key , value := range gameState.PlayerProgress{
			fmt.Println(key , value)
		}
		messageBytes, _ := json.Marshal(startMessage)
		gs.broadcastToRoom(room , messageBytes)
		log.Println("after broadcast for ###startNewGame%%%%")
		fmt.Println(gs.rooms[room], "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		for key , value := range gs.rooms[room]{
			fmt.Println(key , "MTF key")
			fmt.Println(value , "MTF value")
		}
		fmt.Println(gs.games[room].InGameUsers, "{}{}{}{}{}{}")

}


func (gs *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request){
	conn, err := gs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	var username string
	protocols := websocket.Subprotocols(r)
	for _, protocol := range protocols {
		if strings.HasPrefix(protocol, "auth_token:"){
			token := strings.TrimPrefix(protocol , "auth_token:")
			auth, err := verifyTokenWithFastAPI(token)
			fmt.Println("GONNA SEE SOME ERRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR??")
			if err != nil{
				fmt.Println("here is the VALLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLIDDDDDDDDDD")
				fmt.Println(auth.Valid)
				username = auth.UserName
			}
		}else if strings.HasPrefix(protocol, "nickname:"){
			if username == ""{
				username = strings.TrimPrefix(protocol, "nickname:")
			}
		}
	}
	username = fmt.Sprintf("Guest_%d", time.Now().UnixNano()%10000)
	log.Printf("User comes as %s ", username)
	client := &Client{
		conn: conn,
		id: fmt.Sprintf("%s_%d", username, time.Now().UnixNano()),
		username: username,
		room: r.URL.Query().Get("room"),
		sendChan: make(chan []byte, 256),
	}

	joinMessage := struct {
		Type string `json:"type"`
		Username string `json:"username"`
		Content string `json:"content"`
	}{
		Type: "join",
		Username: username,
		Content: fmt.Sprintf("%s joined the game", client.username), 

	}
	
	joinMessageBytes, _ := json.Marshal(joinMessage)
	gs.broadcastToRoom(client.room, joinMessageBytes)



	go gs.readPump(client)
	go gs.writePump(client)

}

func (gs *GameServer) readPump(client *Client){
	defer func(){
		gs.unregister <- client
		client.conn.Close()
	}()

	for {
		messageType, message , err := client.conn.ReadMessage()	
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		if messageType == websocket.TextMessage {
			gs.handleGameMessage(client , message)
		}

	}
}


func (gs *GameServer) writePump(client *Client){
	defer client.conn.Close()
	for {
		select {
		case message, ok := <- client.sendChan:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("error writing message:", err)
				return
			}
			log.Println(string(message), "print in writePum")
		}
	}
}


func (gs *GameServer) userCred(client *Client, messageContent json.RawMessage){
	var result map[string]string
	err := json.Unmarshal(messageContent, &result)
	if err != nil{
		fmt.Println("error while parsing message in wordComplete method")
	}
	client.username = result["username"]
}

func (gs *GameServer) handleGameMessage(client *Client, message []byte){
	var gameMessage struct {
		Type string `json:"type"`
		Content json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(message, &gameMessage) ; err != nil {
		log.Printf("Error parsing message")
		return 
	}
	switch gameMessage.Type{
	case "roomsStatus" :
		gs.roomsStatus(client)
	case "join":
		log.Println("a joining is happening !!!")
		gs.joinPlayer(client , gameMessage.Content)
	case "ready":
		log.Println("player clicked the ready button!!!!")
		gs.readyPlayer(client)
	case "startGame":
		log.Println("New Game!!!!")
		gs.readyToStart(client.room)
	case "wordComplete":
		log.Println("new words come in ! ( print in handleGameMessage)")
		gs.wordComplete(client , gameMessage.Content)
	case "endGame":
		log.Println("+++++++++++++ENG GAME SOLDIER+++++++++++++")
	case "roomStatus":
		log.Println("going into roomStatus Method ( in handleGameMessage)")
		gs.roomStatus(client)
	case "usercred":
		gs.userCred(client , gameMessage.Content)
		
	//case "endGame":
	//	log.Println("End Game!!!!")
	//	gs.endGameMessageHandler()	

	}	


	//enrichedMessage := struct {
	//	Type     string          `json:"type"`
	//	Progress     map[string]*GameState           `json:"progress"`
	//	Username string          `json:"username"`
	//	Content  json.RawMessage `json:"content"`
	//}{
	//	Type: gameMessage.Type,
	//	Progress: gs.games,
	//	Username: client.username,
	//	Content: gameMessage.Content,
	//}
	//enrichedMessageBytes, _ := json.Marshal(enrichedMessage)
	//gs.broadcastToRoom(client.room , enrichedMessageBytes)

	
}


func (gs *GameServer) broadcastToRoom(room string, message []byte){
	log.Println("start of broadcast")
	gs.mutex.Lock()
	log.Println("after Lock")
	defer gs.mutex.Unlock()
	if clients, ok := gs.rooms[room]; ok {
		fmt.Println("found a client!!!!!!!!!!!!!!", clients)
		for _, client := range clients {
			select{
			case client.sendChan <- message:
				var enrichedMessage struct {
						Type     string          `json:"type"`
						Progress     map[string]*GameState           `json:"progress"`
						Username string          `json:"username"`
						Content  json.RawMessage `json:"content"`
					}
				fmt.Println(client.sendChan, "How")
				if err := json.Unmarshal(message , &enrichedMessage) ; err != nil {
					log.Printf("Error parsing message")
					return 
				}
				fmt.Println(enrichedMessage)

			default:
				close(client.sendChan)
				delete(gs.rooms[room], client.id)
				delete(gs.clients, client.id)
			}
		}
	}
	log.Println("no clients in the room")
	fmt.Println(gs.rooms[room] , " the clients are")
	fmt.Println(room, " room is")
}






type AuthRequest struct {
	Scheme string `json:"scheme"`
	Credentials string `json:"credentials"`
}

type AuthResponse struct {
	Valid bool
	UserName string
}

func verifyTokenWithFastAPI(token string) (*AuthResponse, error) {
	authReq := AuthRequest{
		Scheme: "bearer",
		Credentials: token,
	}
	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(
		"http://127.0.0.1:8000/verify",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil , err
	}
	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,err
	}
	return &result, nil
}


func HealthCheck(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w, "Server is okay \n")
}


func main() {


	godotenv.Load(".env")
	portString := os.Getenv("PORT")
    gameServer := NewGameServer("ws://127.0.0.1:8000/ws")

    go gameServer.Run()

    http.HandleFunc("/ws", gameServer.HandleWebSocket) // passing HandleWebSocket method for HandleFunc method ass a value ( that first citizen function kind of things )
	log.Printf("Server starting on port %v", portString)
    log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}


