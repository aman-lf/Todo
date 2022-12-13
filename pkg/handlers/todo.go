package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aman-lf/todo/pkg/models"
)

func ViewTodo(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}
	results := models.GetAllTodos()

	files := []string{
		"pkg/templates/base.gohtml",
		"pkg/templates/index.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", results)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

func CreateTodo(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "POST" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}

	if request.Method == "GET" {
		files := []string{
			"pkg/templates//base.gohtml",
			"pkg/templates//add.gohtml",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Print(err.Error())
			http.Error(response, "Internal Server Error", 500)
			return
		}

		err = ts.ExecuteTemplate(response, "base", nil)
		if err != nil {
			log.Print(err.Error())
			http.Error(response, "Internal Server Error", 500)
		}
	} else if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			log.Print(err.Error())
			return
		}
		item := request.FormValue("item")
		completed_form := request.FormValue("completed")
		var completed bool
		if completed_form == "true" {
			completed = true
		} else {
			completed = false
		}

		models.CreateTodo(item, completed)
		fmt.Fprintf(response, "Data inserted successfully")
	}
}

func UpdateTodo(response http.ResponseWriter, request *http.Request) {
	if request.Method != "PUT" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}
	id := request.URL.Query().Get("id")
	completed_form := request.FormValue("completed")

	var completed bool
	if completed_form == "true" {
		completed = true
	} else {
		completed = false
	}

	models.UpdateTodo(id, completed)
	fmt.Fprintf(response, "Data updated successfully")
}

func UpdateAllTodo(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}

	models.UpdateAllTodo()
	fmt.Fprintf(response, "Data updated successfully")
}

func DeleteTodo(response http.ResponseWriter, request *http.Request) {
	if request.Method != "DELETE" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}
	id := request.URL.Query().Get("id")

	models.DeleteTodo(id)
	fmt.Fprintf(response, "Data deleted successfully")
}
