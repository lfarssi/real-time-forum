package websockets

import (
	"encoding/json"
	"net/http"
	"time"

	"real_time_forum/backend/models"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var userConnections = make(map[int]*websocket.Conn)

func MessageWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		conn.WriteJSON(map[string]any{
			"message": "Something wrong with the message",
			"status":  http.StatusBadRequest,
		})
		return
	}
	userID, ok := r.Context().Value("userId").(int)
	if !ok {
		conn.WriteJSON(map[string]any{
			"message": "You don't have authorization",
			"status":  http.StatusBadRequest,
		})
		return
	}
	userConnections[userID] = conn

	defer func() {
		delete(userConnections, userID)
		conn.Close()
		broadcastStatus(userID, false)
	}()

	broadcastStatus(userID, true)
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			conn.WriteJSON(map[string]any{
				"message": "Cannot Convert Message",
				"status":  http.StatusBadRequest,
			})
			return
		}

		message.SenderID = userID

		switch message.Type {
		case "addMessage":
			err = models.AddMessage(&message)
			if err != nil {
				conn.WriteJSON(map[string]any{
					"message": "Cannot Send Message",
					"status":  http.StatusInternalServerError,
				})
				return
			}
			message.SentAt = time.Now().Format(time.TimeOnly)
			if con, ok := userConnections[message.RecipientID]; ok {
				con.WriteJSON(map[string]any{
					"message":  "Messages Loaded",
					"type":     "newMessage",
					"status":   http.StatusOK,
					"data":     message,
					"isSender": false,
				})
			}
			conn.WriteJSON(map[string]any{
				"message":  "Message Sent",
				"type":     "newMessage",
				"status":   http.StatusOK,
				"data":     message,
				"isSender": true,
			})

		case "loadMessage":
			messages, err := models.GetMessage(message.SenderID, message.RecipientID)
			if err != nil {
				conn.WriteJSON(map[string]any{
					"message": "Error Getting Messages",
					"status":  http.StatusInternalServerError,
				})
				return
			}
			conn.WriteJSON(map[string]any{
				"message": "Messages Loaded",
				"type":    "allMessages",
				"status":  http.StatusOK,
				"data":    messages,
			})
		}

	}
}

func broadcastStatus(userID int, isOnline bool) {
	statusMessage := map[string]any{
		"type":     "userStatus",
		"userID":   userID,
		"isOnline": isOnline,
	}

	for key := range userConnections {
		if conn, ok := userConnections[key]; ok {
			conn.WriteJSON(statusMessage)
		}
	}
}

func OnlineFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value("userId").(int)

	friends, err := models.Friends(userID)
	if err != nil {
		http.Error(w, "Cannot get friends", http.StatusInternalServerError)
		return
	}

	for i, friend := range friends {
		_, friends[i].IsOnline = userConnections[friend.ID]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Friends list",
		"status":  http.StatusOK,
		"data":    friends,
	})
}
