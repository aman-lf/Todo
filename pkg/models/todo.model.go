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

func CreateTodo(item string, completed bool) error {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	todo := bson.D{{"item", item}, {"completed", completed}}
	_, err := Todo.InsertOne(ctx, todo)

	return err

}

func UpdateTodo(id string, completed bool) error {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	updatedTodo := bson.M{"$set": bson.M{"completed": completed}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc bson.D
	updateErr := Todo.FindOneAndUpdate(ctx, filter, updatedTodo, opts).Decode(&updatedDoc)

	return updateErr
}

func UpdateAllTodo() error {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	filter := bson.M{}
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err := Todo.UpdateMany(ctx, filter, update)

	return err
}

func DeleteTodo(id string) error {
	var Todo = db.Database.Collection("todos")
	var ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	_, delete_err := Todo.DeleteOne(ctx, filter)

	return delete_err
}
