package yjs_server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

const (
	messageSync      = 0
	messageAwareness = 1
)

type Document struct {
	Name        string
	Connections map[*websocket.Conn]bool
	Awareness   map[uint32]interface{}
	Content     []byte
	Dirty       bool
	Mutex       sync.Mutex
}

type Server struct {
	Documents map[string]*Document
	Upgrader  websocket.Upgrader
	Mutex     sync.Mutex
}

func InitYjsServer() *Server {
	server := &Server{
		Documents: make(map[string]*Document),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	return server
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	docName := r.URL.Path[1:]
	if docName == "" {
		log.Println("No document name provided")
		return
	}

	doc := s.getOrCreateDocument(docName)

	doc.Mutex.Lock()
	doc.Connections[conn] = true
	doc.Mutex.Unlock()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		if messageType != websocket.BinaryMessage {
			fmt.Println("Got a non-binary message")
			continue
		}

		err = s.handleYjsMessage(conn, doc, message)
		if err != nil {
			log.Println("Yjs message error:", err)
			break
		}
	}

	doc.Mutex.Lock()
	delete(doc.Connections, conn)
	doc.Mutex.Unlock()
}

func (s *Server) getOrCreateDocument(name string) *Document {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if doc, exists := s.Documents[name]; exists {
		return doc
	}

	doc := &Document{
		Name:        name,
		Connections: make(map[*websocket.Conn]bool),
		Awareness:   make(map[uint32]interface{}),
	}
	s.Documents[name] = doc
	return doc
}

func (s *Server) handleYjsMessage(conn *websocket.Conn, doc *Document, message []byte) error {
	if len(message) < 1 {
		return errors.New("empty message")
	}

	messageType := message[0]
	switch messageType {
	case messageSync:
		return s.handleSyncMessage(conn, doc, message[1:])
	case messageAwareness:
		return s.handleAwarenessMessage(conn, doc, message[1:])
	default:
		return errors.New("unknown message type")
	}
}

func (s *Server) handleSyncMessage(conn *websocket.Conn, doc *Document, message []byte) error {
	doc.Mutex.Lock()
	defer doc.Mutex.Unlock()

	for client := range doc.Connections {
		if client != conn {
			_ = client.WriteMessage(websocket.BinaryMessage, append([]byte{messageSync}, message...))
		}
	}

	if !bytes.Equal(doc.Content, message) {
		doc.Content = message
		doc.Dirty = true
	}

	return nil
}

func (s *Server) handleAwarenessMessage(conn *websocket.Conn, doc *Document, message []byte) error {
	if len(message) < 4 {
		return errors.New("invalid awareness message")
	}

	clientID := uint32(message[0])<<24 | uint32(message[1])<<16 | uint32(message[2])<<8 | uint32(message[3])
	awarenessUpdate := message[4:]

	doc.Mutex.Lock()
	doc.Awareness[clientID] = awarenessUpdate
	doc.Mutex.Unlock()

	for client := range doc.Connections {
		if client != conn {
			_ = client.WriteMessage(websocket.BinaryMessage, append([]byte{messageAwareness}, message...))
		}
	}
	return nil
}
