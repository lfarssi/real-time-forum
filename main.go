package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"real_time_forum/backend/database"
	"real_time_forum/backend/router"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		database.DB.Close()
		fmt.Println("\nServer closed succesfully")
		os.Exit(0)
	}()

	router.WebRouter()
	router.ApiRouter()
	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
