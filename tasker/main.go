package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gookit/color.v1"
	"log"
	"os"
	"time"
)

var (
	collection *mongo.Collection
	ctx        = context.TODO()
)

func init() {
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tasker").Collection("tasks")

}

func main() {
	app := &cli.App{
		Name:  "tasker",
		Usage: "sample cli for manage your tasks",
		Commands: cli.Commands{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a new task to storage",
				Action: func(c *cli.Context) error {
					str := c.Args().First()
					if str == "" {
						return errors.New("cannot add empty task")
					}
					task := &Task{
						ID:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdateAt:  time.Now(),
						Text:      str,
						Complete:  false,
					}

					return createTask(task)
				},
			}, {
				Name:    "all",
				Usage:   "list of all tasks",
				Aliases: []string{"l"},
				Action: func(c *cli.Context) error {
					tasks, err := getAll()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Println("non record founded , please add task first")
							return nil
						}

						return err
					}

					printTasks(tasks)
					return nil
				},
			}, {
				Name:    "done",
				Usage:   "Change status of task to completed",
				Aliases: []string{"d"},
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					err := completeTask(text)

					if err != nil {
						return err
					}

					return nil
				},
			}, {
				Name:    "remove",
				Usage:   "delete a task",
				Aliases: []string{"rm"},
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					err := deleteTask(text)

					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func printTasks(tasks []*Task) {
	for i, v := range tasks {
		if v.Complete {
			color.Green.Printf("%d:%s\n", i+1, v.Text)
		} else {
			color.Red.Printf("%d:%s\n", i+1, v.Text)
		}
	}
}
