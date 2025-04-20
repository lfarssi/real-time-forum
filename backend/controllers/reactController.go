package controllers

import (
	"net/http"
	"strconv"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func ReactPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	var react models.React

	var err error
	iduser := r.Context().Value("userId").(int)
	react.UserID = iduser
	react.Status = r.URL.Query().Get("status")
	react.Sender =r.URL.Query().Get("sender")
	if react.Sender == "post" {
		postID, err := strconv.Atoi(r.URL.Query().Get("postID"))
		if err != nil {
			utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": "post id not an integer",
				"status":  http.StatusBadRequest,
			})
			return
		}
		react.PostID = postID

		err = models.InsertReactPost(react)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Cannot react Try Again",
				"status":  http.StatusInternalServerError,
			})
			return
		}

	} else if react.Sender == "comment" {
		react.CommentID,err = strconv.Atoi(r.URL.Query().Get("commentID"))
		if err != nil {
			utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": "comment id not an integer",
				"status":  http.StatusBadRequest,
			})
			return
		}
		err = models.InsertReactComment(react)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Cannot react Try Again",
				"status":  http.StatusInternalServerError,
			})
			return
		}

	}
}
