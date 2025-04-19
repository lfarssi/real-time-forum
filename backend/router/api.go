package router

import (
	"net/http"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
)

func ApiRouter() {
	http.HandleFunc("/api/login", controllers.Login)
	http.HandleFunc("/api/register", controllers.Register)

	http.HandleFunc("/api/addPost", middleware.Authorization(http.HandlerFunc(controllers.AddPost)))
}
