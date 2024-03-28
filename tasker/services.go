package main

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)

	return err
}

func getAll() ([]*Task, error) {
	filter := bson.D{{}}

	return filterTasks(filter)
}

func completeTask(text string) error {
	filter := bson.D{primitive.E{
		Key:   "text",
		Value: text,
	}}

	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "complete", Value: true},
	}}}

	var t Task

	return collection.FindOneAndUpdate(ctx, filter, update).Decode(&t)
}

func deleteTask(text string) error {
	filter := bson.D{bson.E{Key: "text", Value: text}}

	res, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("non tasks deleted")
	}

	return nil

}

func filterTasks(filter interface{}) ([]*Task, error) {
	var tasks []*Task

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return tasks, err
	}

	for cur.Next(ctx) {
		var t Task
		err = cur.Decode(&t)

		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &t)
	}

	if err = cur.Err(); err != nil {
		return tasks, err
	}

	//when cursor exhausted,close that
	err = cur.Close(ctx)
	if err != nil {
		return tasks, err
	}

	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}

	return tasks, nil

}
