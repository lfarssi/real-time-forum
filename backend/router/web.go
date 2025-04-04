package router 

import (

	"net/http"
)

func WebRouter()  {
	http.HandleFunc("/",controller)
}