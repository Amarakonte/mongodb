package config

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func AddParticipant(id_event string, data map[string]interface{}) {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}

	participant := Participants{
		Event:    Event{Id: id_event},
		User:     User{Id: GetUserID(db, data["username"].(string))}, // Passer db en tant que valeur
		Accepted: false,
	}

	_, err = db.Database.Collection("participants").InsertOne(context.Background(), participant)
	if err != nil {
		panic(err)
	}
}

func GetParticipants(db *db, id_event string, data map[string]interface{}) {
	cursor, err := db.Database.Collection("participants").Find(context.Background(), bson.M{"event.id": id_event})
	if err != nil {
		panic(err)
	}

	var participants []Participants
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var participant Participants
		if err := cursor.Decode(&participant); err != nil {
			panic(err)
		}
		participants = append(participants, participant)
	}

	data["participants"] = participants
}

func AcceptParticipation(db *db, id_event string, id_user string) {
	_, err := db.Database.Collection("participants").UpdateOne(context.Background(), bson.M{"event.id": id_event, "user.id": id_user}, bson.M{"$set": bson.M{"accepted": true}})
	if err != nil {
		panic(err)
	}
}

func RemoveParticipant(db *db, id_event string, id_user string) {
	_, err := db.Database.Collection("participants").DeleteOne(context.Background(), bson.M{"event.id": id_event, "user.id": id_user})
	if err != nil {
		panic(err)
	}
}

func IsParticipant(db *db, id_event string, data map[string]interface{}) bool {
	cursor, err := db.Database.Collection("participants").Find(context.Background(), bson.M{"event.id": id_event, "user.id": GetUserID(db, data["username"].(string)), "accepted": true})
	if err != nil {
		panic(err)
	}

	return cursor.Next(context.Background())
}

func HasRequestedParticipation(db *db, id_event string, data map[string]interface{}) bool {
	cursor, err := db.Database.Collection("participants").Find(context.Background(), bson.M{"event.id": id_event, "user.id": GetUserID(db, data["username"].(string))})
	if err != nil {
		panic(err)
	}

	return cursor.Next(context.Background())
}
