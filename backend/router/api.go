package router

import (
	"net/http"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
)

func ApiRouter() {
	http.HandleFunc("/api/login", controllers.LoginController)
	http.HandleFunc("/api/register", controllers.RegisterController)
	http.HandleFunc("/api/isLogged", middleware.IsLogged)

	http.HandleFunc("/api/addPost", middleware.Authorization(http.HandlerFunc(controllers.AddPostController)))
	http.HandleFunc("/api/getPosts", middleware.Authorization(http.HandlerFunc(controllers.GetPostController)))
	http.HandleFunc("/api/getLikedPosts", middleware.Authorization(http.HandlerFunc(controllers.GetLikedPostController)))
	http.HandleFunc("/api/getCreatedPosts", middleware.Authorization(http.HandlerFunc(controllers.GetCreatedPostController)))
	http.HandleFunc("/api/getPostsByCategory", middleware.Authorization(http.HandlerFunc(controllers.GetPostByCategoryController)))

	http.HandleFunc("/api/addComment", middleware.Authorization(http.HandlerFunc(controllers.AddCommentController)))
	http.HandleFunc("/api/getComments", middleware.Authorization(http.HandlerFunc(controllers.GetCommnetsController)))

	http.HandleFunc("/api/addLike", middleware.Authorization(http.HandlerFunc(controllers.ReactPostController)))
}
