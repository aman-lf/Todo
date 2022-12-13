package models

import (
	"context"
	"log"
	"time"

	"github.com/aman-lf/todo/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllTodos() []interface{} {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	cursor, err := Todo.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
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

		// Converting ObjectID to string
		elem["_id"] = elem["_id"].(primitive.ObjectID).Hex()
		results = append(results, elem)
	}

	return results
}

func CreateTodo(item string, completed bool) {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	todo := bson.D{{"item", item}, {"completed", completed}}
	_, err := Todo.InsertOne(ctx, todo)
	if err != nil {
		log.Print(err.Error())
		return
	}
}

func UpdateTodo(id string, completed bool) {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"_id": objID}
	updatedTodo := bson.M{"$set": bson.M{"completed": completed}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc bson.D
	updateErr := Todo.FindOneAndUpdate(ctx, filter, updatedTodo, opts).Decode(&updatedDoc)
	if updateErr != nil {
		panic(updateErr)
	}
}

func UpdateAllTodo() {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	filter := bson.M{}
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err := Todo.UpdateMany(ctx, filter, update)
	if err != nil {
		panic(err)
	}
}

func DeleteTodo(id string) {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"_id": objID}

	_, delete_err := Todo.DeleteOne(ctx, filter)
	if delete_err != nil {
		panic(delete_err)
	}
}
