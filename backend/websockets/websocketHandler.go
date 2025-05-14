package websockets

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"real_time_forum/backend/models"

	"github.com/gorilla/websocket"
)

var (
	upgrade         = websocket.Upgrader{}
	userConnections = make(map[int][]*websocket.Conn)
	mu              sync.Mutex
)

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

	mu.Lock()
	userConnections[userID] = append(userConnections[userID], conn)
	mu.Unlock()

	broadcastStatus(userID, true)

	unreadCounts, err := models.GetUnreadCountsPerFriend2(userID)
	if err == nil {
		for _, conn := range userConnections[userID] {
			conn.WriteJSON(map[string]any{
				"type":   "unreadCounts",
				"counts": unreadCounts,
			})
		}
	}
	defer removeConnection(userID, conn)

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

		if message.LastID == -1 {
			message.LastID, err = models.GetLastMessageID()
			if err != nil {
				conn.WriteJSON(map[string]any{
					"message": err.Error(),
					"status":  http.StatusInternalServerError,
				})
				return
			}
		}

		// message.SenderID = userID
		message.Content = html.EscapeString(message.Content)

		switch message.Type {
		case "addMessage":
			if strings.TrimSpace(message.Content) == "" {
				conn.WriteJSON(map[string]any{
					"message": "Cannot Send Empty Message",
					"status":  http.StatusBadRequest,
					"type":    "errMessage",
				})
			} else if len(strings.TrimSpace(message.Content)) > 1000 {
				conn.WriteJSON(map[string]any{
					"message": "Cannot Send Over 1000 Character",
					"status":  http.StatusBadRequest,
					"type":    "errMessage",
				})
			} else {
				err = models.AddMessage(&message)
				if err != nil {
					conn.WriteJSON(map[string]any{
						"message": "Cannot Send Message",
						"status":  http.StatusInternalServerError,
						"type":    "errMessage",
					})
					return
				}
				message.SentAt = time.Now()
				unreadCounts, _ := models.GetUnreadCountsPerFriend(message.RecipientID, userID)

				if conns, ok := userConnections[message.RecipientID]; ok {
					for _, c := range conns {
						c.WriteJSON(map[string]any{
							"message":  "Messages Loaded",
							"type":     "newMessage",
							"status":   http.StatusOK,
							"data":     message,
							"counts":   unreadCounts,
							"isSender": false,
						})
					}
				}
				for _, c := range userConnections[message.SenderID] {
					c.WriteJSON(map[string]any{
						"message":  "Message Sent",
						"type":     "newMessage",
						"status":   http.StatusOK,
						"data":     message,
						"counts":   unreadCounts,
						"isSender": true,
					})
				}
			}

		case "loadMessage":

			messages, err := models.GetMessage(userID, message.RecipientID, message.LastID)
			if err != nil {
				for _, c := range userConnections[userID] {
					c.WriteJSON(map[string]any{
						"message": "Error Getting Messages",
						"status":  http.StatusInternalServerError,
						"type":    "errMessage",
					})
				}
				return
			}

			conn.WriteJSON(map[string]any{
				"message": "Messages Loaded",
				"type":    "allMessages",
				"status":  http.StatusOK,
				"data":    messages,
			})
		case "Typing":
			for _, c := range userConnections[message.RecipientID] {
				c.WriteJSON(map[string]any{
					"message": "is Typing",
					"status":  http.StatusOK,
					"userId":  message.SenderID,
					"type":    "isTyping",
				})
			}
		case "StopTyping":
			for _, c := range userConnections[message.RecipientID] {
				c.WriteJSON(map[string]any{
					"message": "is Not Typing",
					"status":  http.StatusOK,
					"userId":  message.SenderID,
					"type":    "isNotTyping",
				})
			}
		case "updateMessage":
			err := models.UpdateMessage(message.SenderID, message.RecipientID, message.Status)
			if err != nil {
				conn.WriteJSON(map[string]any{
					"message": "Messages Updated",
					"type":    "messageUpdated",
					"status":  http.StatusOK,
				})
			}
		case "logout":
			conn.WriteJSON(map[string]any{
				"message": "user logged out",
				"type":    "refreshFriends",
				"status":  http.StatusOK,
			})

		}

	}
}

func removeConnection(userID int, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	conns := userConnections[userID]
	for i, c := range conns {
		if c == conn {
			userConnections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(userConnections[userID]) == 0 {
		delete(userConnections, userID)
		broadcastStatus(userID, false)
	}
}

func broadcastStatus(userID int, isOnline bool) {
	// Get the user info for the user whose status changed
	user, err := models.GetUserByID(userID)
	if err != nil {
		// handle error, maybe log and return
		return
	}

	friends, _ := models.Friends(userID)

	statusMessage := map[string]any{
		"type":      "userStatus",
		"userID":    userID,
		"isOnline":  isOnline,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"userName":  user.UserName,
	}

	for _, friend := range friends {
		if conns, ok := userConnections[friend.ID]; ok {
			for _, c := range conns {
				c.WriteJSON(statusMessage)
			}
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
