package router

import (
	"net/http"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
	"real_time_forum/backend/websockets"
)

func ApiRouter() {
	http.HandleFunc("/api/login", controllers.LoginController)
	http.HandleFunc("/api/register", controllers.RegisterController)
	http.HandleFunc("/api/logout", controllers.LogoutController)
	http.HandleFunc("/api/isLogged", middleware.IsLogged)

	http.HandleFunc("/api/addPost", middleware.Authorization(http.HandlerFunc(controllers.AddPostController)))
	http.HandleFunc("/api/getPosts", middleware.Authorization(http.HandlerFunc(controllers.GetPostController)))
	http.HandleFunc("/api/getLikedPosts", middleware.Authorization(http.HandlerFunc(controllers.GetLikedPostController)))
	http.HandleFunc("/api/getCreatedPosts", middleware.Authorization(http.HandlerFunc(controllers.GetCreatedPostController)))
	http.HandleFunc("/api/getPostsByCategory", middleware.Authorization(http.HandlerFunc(controllers.GetPostByCategoryController)))

	http.HandleFunc("/api/addComment", middleware.Authorization(middleware.RateLimit(http.HandlerFunc(controllers.AddCommentController))))
	http.HandleFunc("/api/getComments", middleware.Authorization(http.HandlerFunc(controllers.GetCommnetsController)))
	http.HandleFunc("/api/getCategories", middleware.Authorization(http.HandlerFunc(controllers.CategoryController)))

	http.HandleFunc("/api/addLike", middleware.Authorization(http.HandlerFunc(controllers.ReactPostController)))

	http.HandleFunc("/api/getFriends", middleware.Authorization(http.HandlerFunc(websockets.OnlineFriends)))

	http.HandleFunc("/ws/messages", middleware.Authorization(http.HandlerFunc(websockets.MessageWebSocketHandler)))
}
