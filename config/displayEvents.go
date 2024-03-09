package config

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// DisplayEvents récupère les événements de la base de données et les affiche en fonction du critère de recherche
func DisplayEvents(data map[string]interface{}, searching string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var events []Event

	// Définir le filtre de recherche en fonction du critère de recherche
	filter := bson.M{}
	if searching != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": searching, "$options": "i"}},
			{"date_start": bson.M{"$regex": searching, "$options": "i"}},
		}
	}

	// Récupérer les événements correspondant au filtre
	cursor, err := db.Database.Collection("event").Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	// Parcourir les événements récupérés et les ajouter à la liste des événements à afficher
	for cursor.Next(context.Background()) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return err
		}

		// Remplacer les sauts de ligne par <br> pour l'affichage HTML
		event.Description = strings.Replace(strings.Replace(event.Description, "\r", "", -1), "\n", "<br>", -1)

		events = append(events, event)
	}

	// Ajouter la liste des événements à afficher aux données à envoyer au template HTML
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
