package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Init() error {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+1.8.0")

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return fmt.Errorf("connect mongodb error: %s", err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return fmt.Errorf("ping mongodb error: %s", err)
	}

	MongoClient = client
	return nil
}

func Close() error {
	if err := MongoClient.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("disconnect mongodb error: %s", err)
	}
	return nil
}

func Query(collection *mongo.Collection, filter interface{}) (*mongo.Cursor, error) {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("%s find error: %s", collection.Name(), err)
	}

	return cursor, nil
}

func Insert(collection *mongo.Collection, documents []interface{}) error {
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		return fmt.Errorf("%s insert error: %s", collection.Name(), err)
	}

	return nil
}

func Delete(collection *mongo.Collection, filter interface{}) (int64, error) {
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return -1, fmt.Errorf("%s delete error: %s", collection.Name(), err)
	}

	return result.DeletedCount, nil
}

func UpdateOne(collcetion *mongo.Collection, filter bson.D, update bson.D) error {
	_, err := collcetion.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("%s update error: %s", collcetion.Name(), err)
	}

	return nil
}

func ListenChangeStream(collection *mongo.Collection, filter interface{}) error {
	_, err := collection.Watch(context.Background(), mongo.Pipeline{{{Key: "$match", Value: filter}}},
		options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		return err
	}
	return nil
}
