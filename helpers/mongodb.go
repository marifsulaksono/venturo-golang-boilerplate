package helpers

import (
	"context"
	"simple-crud-rnd/config"
	"simple-crud-rnd/structs"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertLogDataToMongoDB(ctx context.Context, db, collection string, data *structs.LogEntry) error {
	_, err := config.MongoCLI.Database(db).Collection(collection).InsertOne(ctx, bson.D{
		{Key: "url", Value: data.URL},
		{Key: "path", Value: data.Path},
		{Key: "ip", Value: data.IP},
		{Key: "users", Value: data.User},
		{Key: "body", Value: data.Body},
		{Key: "response", Value: data.Response},
	})

	return err
}
