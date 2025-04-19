package router

import (
	"net/http"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
)

func ApiRouter() {
	http.HandleFunc("/api/login", controllers.LoginController)
	http.HandleFunc("/api/register", controllers.RegisterController)

	http.HandleFunc("/api/addPost", middleware.Authorization(http.HandlerFunc(controllers.AddPostController)))
}
