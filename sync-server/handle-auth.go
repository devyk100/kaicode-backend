package sync_server

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("your-very-long-random-secret")

func ParseToken(tokenString string) (int64, error) {
	// Parse the token and validate signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure HS256 signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SECRET_KEY, nil
	})
	if err != nil {
		return -1, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims as map
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["userId"].(float64); ok {
			return int64(userIDFloat), nil
		}
		return -1, fmt.Errorf("userId not found or invalid in token")
	}

	return -1, fmt.Errorf("invalid token")
}

func HandleWSAuth(conn *SecureWSConn) (int64, string, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Client disconnected")
		return -1, "", err
	}
	var ReceivedAuthMessage Event
	err = json.Unmarshal(msg, &ReceivedAuthMessage)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return -1, "", err
	}
	if ReceivedAuthMessage.Event == "auth" {
		fmt.Println("Received auth message")
	} else {
		return 0, "", fmt.Errorf("invalid event")
	}
	userId, err := ParseToken(ReceivedAuthMessage.Token)
	if err != nil {
		return 0, "", err
	}
	fmt.Println("Was authorized using", userId, "user", userId)
	return userId, ReceivedAuthMessage.RoomName, nil
}
