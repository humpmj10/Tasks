package main

import (
	"Tasks/db"
	"Tasks/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/tasks", handlers.GetTasksHandler).Methods(http.MethodGet)
	router.HandleFunc("/task", handlers.CreateTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/task/{taskID}", handlers.DeleteTaskHandler).Methods(http.MethodDelete)

	n := negroni.Classic()
	n.UseHandler(router)

	err := http.ListenAndServe(":8080", n)
	if err != nil {
		return
	}
	fmt.Println("Server running on port 8080")
}
