package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aman-lf/todo/pkg/db"
	"github.com/aman-lf/todo/pkg/handlers"
)

func main() {
	db.InitDb()
	http.HandleFunc("/", handlers.ViewTodo)
	http.HandleFunc("/add", handlers.CreateTodo)
	http.HandleFunc("/update/", handlers.UpdateTodo)
	http.HandleFunc("/updateall", handlers.UpdateAllTodo)
	http.HandleFunc("/delete/", handlers.DeleteTodo)

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
