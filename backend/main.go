package main

func main()  {
	database.DatabaseExecution()
	defer database.CloseDatabase()
	router.WebRouter()
	router.ApiRouter()
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err starting the server : ", err)
		return
	}
}
