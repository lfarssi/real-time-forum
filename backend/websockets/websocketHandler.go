package websockets

import (
	"net/http"
	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func MessageWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Something wrong with the message",
			"status":  http.StatusBadRequest,
		})
		return
	}
	defer conn.Close()
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": "Cannot Convert Message",
				"status":  http.StatusBadRequest,
			})
			return
		}

		sender := r.Context().Value("userId").(int)
		message.SenderID=sender
		switch message.Type{
		case "addMessage":
			err = models.AddMessage(&message)
			if err != nil {
				utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
					"message": "Cannot Send Message",
					"status":  http.StatusInternalServerError,
				})
				return
			}
			utils.ResponseJSON(w, http.StatusOK, map[string]any{
				"message": "Message Sent",
				"status":  http.StatusOK,
			})
		case "loadMessage":
			messages, err := models.GetMessage(message.SenderID, message.RecipientID)
			if err!= nil{
				utils.ResponseJSON(w, http.StatusOK, map[string]any{
					"message": "Getting  Messages",
					"status":  http.StatusOK,
				})
				return
			}
			utils.ResponseJSON(w, http.StatusOK, map[string]any{
				"message": "Message Loaded",
				"status":  http.StatusOK,
				"data":    messages,
			})

		}
        
		
	}
}
