package sync_server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"kc-backend/sync-server/sqs"
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
	Input    string `json:"input"`     // input in the case this is a judge event
	RoomName string `json:"room_name"` // at the time of auth itself
}

type Job struct {
	SubmissionId string `json:"submission_id"`
	Code         string `json:"code"`
	Language     string `json:"language"`
	RoomName     string `json:"room_name"`
	Input        string `json:"input"`
	// some aws sqs handle, to delete this from the original SQS later
}

type SecureWSConn struct {
	*websocket.Conn
	Mutex sync.Mutex
}

type SyncServer struct {
	RedisClient           *redis.Client
	Upgrader              websocket.Upgrader
	Mutex                 sync.Mutex
	SubscriptionsMutex    sync.Mutex
	ActiveRooms           map[string]bool
	Subscriptions         map[string][]*SecureWSConn
	DelayedPersistRoutine map[string]*time.Timer
}

func InitSyncServer() *SyncServer {
	sqs.InitSQSClient()
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

	secureConn := &SecureWSConn{conn, sync.Mutex{}}
	_, roomName, err := HandleWSAuth(secureConn)
	s.HandleSubscribe(roomName, secureConn)
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
		case "judge":
			judgeJob := Job{
				SubmissionId: "trash",
				Code:         ReceivedMessage.Content,
				Language:     ReceivedMessage.Language,
				RoomName:     ReceivedMessage.RoomName,
				Input:        ReceivedMessage.Input,
			}
			val, err := json.Marshal(judgeJob)
			if err != nil {
				return
			}
			err = sqs.SendMessage(string(val), roomName)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
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

			//fmt.Println("Received sync message", "from", userId, ReceivedMessage.Content)
			err := s.RedisClient.Set(context.Background(), "doc:"+(roomName), ReceivedMessage.Content, PERSIST_DEBOUNCE_INTERVAL).Err()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			s.Mutex.Unlock()
		}
	}
}
