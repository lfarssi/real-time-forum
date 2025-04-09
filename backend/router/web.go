package router

import (
	"net/http"
	"real_time_forum/backend/controllers"
)

func WebRouter()  {
	http.HandleFunc("/",controllers.HomePageController)
}	