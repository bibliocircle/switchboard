package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func databaseFromConnectionString(connectionString string) string {
	tokens := strings.Split(connectionString, "/")
	return tokens[len(tokens)-1]
}

func Disconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func Connect(ctx context.Context) error {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	} else {
		fmt.Println("successfully connected to the database!")
	}
	Database = Client.Database(databaseFromConnectionString(connectionString))
	return nil
}
