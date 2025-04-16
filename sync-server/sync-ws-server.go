package sync_server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"sync"
	"time"
)

type Event struct {
	Event    string `json:"event"`     // sync, language, judge, auth
	Content  string `json:"content"`   // if sync
	Language string `json:"language"`  // if the language update
	Token    string `json:"token"`     // if token
	RoomName string `json:"room_name"` // at the time of auth itself
}

type SyncServer struct {
	RedisClient           *redis.Client
	Upgrader              websocket.Upgrader
	Mutex                 sync.Mutex
	DelayedPersistRoutine map[string]*time.Timer
}

func InitSyncServer() *SyncServer {
	rdb := redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_URL"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		TLSConfig: &tls.Config{},
		DB:        0,
	})
	syncServer := &SyncServer{
		RedisClient: rdb,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		DelayedPersistRoutine: make(map[string]*time.Timer),
	}
	return syncServer
}

var PERSIST_DEBOUNCE_INTERVAL = 10 * time.Minute

func (s *SyncServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	userId, err := HandleWSAuth(conn)
	if err != nil {
		fmt.Println(err.Error(), "not authorized")
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if messageType == websocket.BinaryMessage {
			fmt.Println(string(message), "is a binary message")
			return
		}
		var ReceivedMessage Event
		err = json.Unmarshal(message, &ReceivedMessage)

		switch ReceivedMessage.Event {
		case "sync":
			s.Mutex.Lock()
			roomName := ReceivedMessage.RoomName
			// lazily postpone the DB persistence, until no longer being used
			if timer, exists := s.DelayedPersistRoutine[roomName]; exists {
				timer.Stop()
			}
			s.DelayedPersistRoutine[roomName] = time.AfterFunc(PERSIST_DEBOUNCE_INTERVAL, func() {
				// DB persist logic
				fmt.Println("Performing DB persistence", "for the room", roomName)

				delete(s.DelayedPersistRoutine, roomName)
			})

			fmt.Println("Received sync message", "from", userId, ReceivedMessage.Content)
			err := s.RedisClient.Set(context.Background(), "doc:"+(roomName), ReceivedMessage.Content, PERSIST_DEBOUNCE_INTERVAL).Err()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			s.Mutex.Unlock()
		}
	}
}
