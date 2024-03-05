package config

import (
	"context"
)

func AddEvent(input_title string, input_description string, input_date string, data map[string]interface{}) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	event := Event{
		Title:       input_title,
		Description: input_description,
		Date_start:  input_date,
		// Vous devrez peut-être ajouter le reste des champs de l'événement ici
	}

	// Insérer l'événement dans la collection "event"
	_, err = db.Database.Collection("event").InsertOne(context.Background(), event)
	if err != nil {
		return err
	}

	// Mettre à jour la variable de données pour indiquer que l'événement a été créé
	data["eventCreated"] = true

	return nil
}
