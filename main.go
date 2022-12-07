package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func InitDataLayer() *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to Database")
	}
	return client
}

func getAllTodos() []interface{} {
	client := InitDataLayer()
	coll := client.Database("todoTraining").Collection("todos")

	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []interface{}
	for cursor.Next(ctx) {
		var elem bson.M
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	return results
}

func insertTodo(item string, completed bool) {
	client := InitDataLayer()
	coll := client.Database("todoTraining").Collection("todos")

	todo := bson.D{{"item", item}, {"completed", completed}}
	_, err := coll.InsertOne(ctx, todo)
	if err != nil {
		log.Print(err.Error())
		return
	}
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}
	results := getAllTodos()

	files := []string{
		"static/base.gohtml",
		"static/index.gohtml",
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

func addHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "POST" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}

	if request.Method == "GET" {
		files := []string{
			"static/base.gohtml",
			"static/add.gohtml",
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

		insertTodo(item, completed)
		fmt.Fprintf(response, "Data inserted successfully")
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
