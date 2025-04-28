package sync_server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func (s *SyncServer) KeepSubscriptionAlive(roomName string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("Subscribing to this room", roomName)
	pubsub := s.RedisClient.Subscribe(ctx, "roomName")
	defer pubsub.Close()

	ch := pubsub.Channel()

	fmt.Println("started the subscription goroutine")
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload, "is the message from redis")
		s.SubscriptionsMutex.Lock()

		//alternative to map[*SecureWSConn]bool, more memory efficient
		toRemove := make(map[*SecureWSConn]struct{})

		for _, conn := range s.Subscriptions[roomName] {
			conn.Mutex.Lock()
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				log.Printf("Error sending message to WebSocket: %v", err)
				toRemove[conn] = struct{}{}
			} else {
				conn.Mutex.Unlock()
			}
		}

		var newSubs []*SecureWSConn
		for _, conn := range s.Subscriptions[roomName] {
			if _, found := toRemove[conn]; !found {
				newSubs = append(newSubs, conn) // Keep only the connections that weren't marked for removal
			}
		}

		s.Subscriptions[roomName] = newSubs // hope that garbage collector clears that old shit

		if len(s.Subscriptions[roomName]) == 0 {
			s.ActiveRooms[roomName] = false
			return
		}

		s.SubscriptionsMutex.Unlock()
	}

}

func (s *SyncServer) HandleSubscribe(roomName string, conn *SecureWSConn) {
	s.SubscriptionsMutex.Lock()
	if s.Subscriptions == nil {
		s.Subscriptions = make(map[string][]*SecureWSConn)
	}
	if s.ActiveRooms == nil {
		s.ActiveRooms = make(map[string]bool)
	}

	if !s.ActiveRooms[roomName] {
		s.ActiveRooms[roomName] = true
		go s.KeepSubscriptionAlive(roomName)
	}

	s.Subscriptions[roomName] = append(s.Subscriptions[roomName], conn)
	s.SubscriptionsMutex.Unlock()
}
