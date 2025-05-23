package controllers

import (
	"encoding/json"
	"net/http"

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

	if err = json.NewDecoder(r.Body).Decode(&react); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	react.UserID= r.Context().Value("userId").(int)
	var like, dislike []*models.React
	var clike, cdislike []*models.React
	if react.Sender == "post" {
		err = models.InsertReactPost(react)

		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Cannot react Try Again",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		like,err = models.GetReactionPost(react.PostID,"like")
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error getting number of reaction",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		dislike, err = models.GetReactionPost(react.PostID,"dislike")
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error getting number of reaction",
				"status":  http.StatusInternalServerError,
			})
			return
		}

	} else if react.Sender == "comment" {
		err = models.InsertReactComment(react)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Cannot react Try Again",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		clike,err = models.GetReactionComment(react.CommentID,"like")
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error getting number of reaction",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		cdislike, err = models.GetReactionComment(react.CommentID,"dislike")
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error getting number of reaction",
				"status":  http.StatusInternalServerError,
			})
			return
		}
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "React added successfully!",
		"status":  http.StatusOK,
		"data": map[string]any{
			"nbLikes": len(like),
			"nbDislikes": len(dislike),
			"cnbLikes": len(clike),
			"cnbDislikes": len(cdislike),

	
	},

	})
}
