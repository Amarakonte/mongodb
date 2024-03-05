package config

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func DisplayEvents(data map[string]interface{}, searching string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var events []Event

	filter := bson.M{}
	if searching != "" {
		filter["title"] = searching
	}

	cursor, err := db.Database.Collection("event").Find(context.Background(), filter)
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return err
		}

		// Remplace les \n par des <br> pour sauter des lignes en html
		event.Description = strings.Replace(strings.Replace(event.Description, "\r", "", -1), "\n", "<br>", -1)

		events = append(events, event)
	}

	data["events"] = events

	return nil
}

func GetEvent(data map[string]interface{}, id_event string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var event Event
	err = db.Database.Collection("event").FindOne(context.Background(), bson.M{"_id": id_event}).Decode(&event)
	if err != nil {
		return err
	}

	if event.Title != "" {
		data["event"] = event
	}

	data["Id_event"] = id_event

	return nil
}

func GetComments(db *db, id_event string, data map[string]interface{}) error {
	cursor, err := db.Database.Collection("comment").Find(context.Background(), bson.M{"eventID": id_event})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	var comments []Comment

	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return err
		}

		// Remplace les \n par des <br> pour sauter des lignes en html
		comment.Content = strings.Replace(strings.Replace(comment.Content, "\r", "", -1), "\n", "<br>", -1)

		comments = append(comments, comment)
	}

	data["comments"] = comments

	return nil
}
