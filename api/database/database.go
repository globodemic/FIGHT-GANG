package database

import (
	"context"
	"log"
	"time"

	model "../models"
	response "../models/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

//GetDatabase returns a MongoDB Client
func GetDatabase() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
}

//CreatePlayer function to instert Player to MongoDb
func CreatePlayer(player model.Player) string {
	collection := client.Database("fight-gang").Collection("players")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, player)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	//ID of the object document
	objectID := result.InsertedID.(primitive.ObjectID).Hex()
	return objectID

}

//GetPlayer Get Player's data
func GetPlayer(cred model.Credentials, params map[string]string) model.Player {

	//Get ObjectID from string
	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.D{
		{"_id", id},
		{"name", cred.Name},
		{"password", cred.Password},
	}

	var player model.Player
	collection := client.Database("fight-gang").Collection("players")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resultPlayer := collection.FindOne(ctx, filter)
	resultPlayer.Decode(&player)
	return player

}

func GetActivePlayer(params string) response.ResponsePlayer {

	id, _ := primitive.ObjectIDFromHex(params)
	filter := bson.M{"_id": id}
	var player response.ResponsePlayer
	collection := client.Database("fight-gang").Collection("players")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resultPlayer := collection.FindOne(ctx, filter)
	resultPlayer.Decode(&player)
	return player

}

//GetPlayers get all players
func GetPlayers() []*response.ResponsePlayer {
	filter := bson.M{}
	var players []*response.ResponsePlayer
	collection := client.Database("fight-gang").Collection("players")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resultPlayers, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}

	for resultPlayers.Next(ctx) {
		var player response.ResponsePlayer
		err = resultPlayers.Decode(&player)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		players = append(players, &player)
	}
	return players

}

//ChangePlayerAuth change player's alias or password
func ChangePlayerAuth(cred model.Credentials, params map[string]string, player model.Player) int64 {
	//Get ObjectID from string
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
		return 0

	}
	filter := bson.D{
		{"_id", id},
		{"name", cred.Name},
		{"password", cred.Password},
	}
	update := bson.M{
		"$set": bson.M{
			"password": player.Password,
			"alias":    player.Alias,
		},
	}

	collection := client.Database("fight-gang").Collection("players")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	updateResult, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Fatal(err)
	}
	return updateResult.ModifiedCount
}

//DeletePlayer delete one existing Player
func DeletePlayer(params map[string]string) int64 {

	//Get ObjectID from string
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
		return 0

	}

	filter := bson.M{"_id": id}
	collection := client.Database("fight-gang").Collection("players")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal("Error on deleting Player", err)
	}
	return deleteResult.DeletedCount
}
