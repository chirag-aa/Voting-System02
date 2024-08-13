package vote

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	ctx      = context.Background()
	rdb      *redis.Client
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type VoteMessage struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type Session struct {
	Votes map[string]int
	lock  sync.RWMutex
}

var sessions = make(map[string]*Session)
var sessionsLock sync.RWMutex

func InitializeRedis(redisClient *redis.Client) {
	rdb = redisClient
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while handling connection:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading message:", err)
			break
		}

		var voteMsg VoteMessage
		if err := json.Unmarshal(msg, &voteMsg); err != nil {
			log.Println("Error while unmarshalling message:", err)
			continue
		}

		// Handle vote message
		handleVoteMessage(conn, voteMsg)
	}
}

func handleVoteMessage(conn *websocket.Conn, msg VoteMessage) {
	sessionsLock.RLock()
	session, exists := sessions[msg.Data]
	sessionsLock.RUnlock()
	if !exists {
		session = &Session{Votes: make(map[string]int)}
		sessionsLock.Lock()
		sessions[msg.Data] = session
		sessionsLock.Unlock()
	}

	session.lock.Lock()
	if msg.Action == "vote" {
		session.Votes[msg.Data]++
		broadcastVoteResults(session)
	}
	session.lock.Unlock()
}

func broadcastVoteResults(session *Session) {
	sessionsLock.RLock()
	for _, sess := range sessions {
		sess.lock.RLock()
		// Broadcasting to all connected users in the session
		sess.lock.RUnlock()
	}
	sessionsLock.RUnlock()
}
