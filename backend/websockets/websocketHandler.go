package websockets

import (
	"net/http"

	"real_time_forum/backend/models"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var Clients = make(map[*websocket.Conn]bool)

func MessageWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		conn.WriteJSON(map[string]any{
			"message": "Something wrong with the message",
			"status":  http.StatusBadRequest,
		})
		return
	}
	Clients[conn] = true
	defer func() {
		delete(Clients, conn)
		conn.Close()
	}()
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

		message.SenderID = r.Context().Value("userId").(int)

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
			conn.WriteJSON(map[string]any{
				"message": "Message Sent",
				"type":    "newMessage",
				"status":  http.StatusOK,
				"data":    message,
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
