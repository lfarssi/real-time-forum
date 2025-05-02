package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"html"


	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func GetPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	userID := r.Context().Value("userId").(int)

	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * 10
	posts, err := models.GetPosts(userID, offset)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error Getting Post",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Posts retrieved successfully",
		"status":  http.StatusOK,
		"data":    posts,
	})
}

func GetLikedPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	userID := r.Context().Value("userId").(int)

	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * 10
	posts, err := models.LikedPost(userID, offset)
	if err != nil {

		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error Getting Post",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Posts retrieved successfully",
		"status":  http.StatusOK,
		"data":    posts,
	})
}

func GetCreatedPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	userID := r.Context().Value("userId").(int)

	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * 10
	posts, err := models.CreatedPost(userID, offset)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error Getting Post",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Posts retrieved successfully",
		"status":  http.StatusOK,
		"data":    posts,
	})
}

func GetPostByCategoryController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	r.ParseForm()

	categories := r.Form["categories"]
	postSet := make(map[int]struct{})
	var posts []models.Post

	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * 10
	for _, category := range categories {
		idCategorie, err := strconv.Atoi(category)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error converting categorie id",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		postTemp, err := models.GetPostsByCategory(idCategorie, offset)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error Getting Post",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// Avoid duplicate posts by using a map
		for _, post := range postTemp {
			if _, exist := postSet[post.ID]; !exist {
				posts = append(posts, post)
				postSet[post.ID] = struct{}{}
			}
		}
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].DateCreation > posts[j].DateCreation
	})

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Posts retrieved successfully",
		"status":  http.StatusOK,
		"data":    posts,
	})
}

func AddPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	var post *models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)

	if !verifyPostData(w, post.Title, post.Content, post.Categories) {
		return
	}

	post.UserID = r.Context().Value("userId").(int)
	post.Title= html.EscapeString(post.Title)
	post.Content= html.EscapeString(post.Content)
	err := models.AddPost(post)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post added successfully!",
		"status":  http.StatusOK,
	})
}

func verifyPostData(w http.ResponseWriter, title, content string, categories []string) bool {
	var message models.ValidationMessagesAddPost
	isValid := true

	if title == "" {
		message.TitleMessage = "Title is required."
		isValid = false
	} else if len([]rune(title)) > 100 {
		message.TitleMessage = "Title must be less than or equal to 100 characters."
		isValid = false
	}

	if content == "" {
		message.ContentMessage = "Content is required."
		isValid = false
	} else if len([]rune(content)) > 1000 {
		message.ContentMessage = "Content must be less than or equal to 1000 characters."
		isValid = false
	}

	if len(categories) == 0 {
		message.CategoryMessage = "At least one category is required."
		isValid = false
	} else if len(categories) > 100 {
		message.CategoryMessage = "You can add up to 100 categories only."
		isValid = false
	} else {
		for _, category := range categories {
			categoryID, err := strconv.Atoi(category)
			if err != nil {
				message.CategoryMessage = "Each category must be a valid number."
				isValid = false
				break
			}

			if !models.IsExistsCategory(categoryID) {
				message.CategoryMessage = "One or more selected categories do not exist."
				isValid = false
				break
			}
		}
	}

	if !isValid {
		utils.ResponseJSON(w, http.StatusBadRequest, message)
		return false
	}

	return true
}
