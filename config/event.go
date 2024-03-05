package config

import (
	"context"
	"time"
)

func AddCommentOnEvent(id_event string, input_commentaire string, data_PageEvent map[string]interface{}) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	comment := Comment{
		Content:      input_commentaire,
		CreationDate: time.Now().Format(time.RFC3339),
		// Vous devrez définir l'utilisateur associé au commentaire ici
		// Et spécifier l'ID de l'événement associé au commentaire
	}

	// Insérer le commentaire dans la collection "comment"
	_, err = db.Database.Collection("comment").InsertOne(context.Background(), comment)
	if err != nil {
		return err
	}

	data_PageEvent["commentaireCreated"] = true

	return nil
}

func UpdateNote(db *db, data map[string]interface{}, id_event string, newNote string) error {
	// La logique pour mettre à jour la note de l'événement reste la même, utilisez simplement les fonctions MongoDB
	// pour effectuer la mise à jour des documents dans la collection "event"
	return nil
}

func GetAllEvents(db *db) ([]Event, error) {
	// Utilisez la fonction Find de la collection "event" pour récupérer tous les événements dans MongoDB
	// et décodez les résultats dans une slice d'Event
	return nil, nil
}

func DeleteEvent(db *db, event_id string) error {
	// Utilisez la fonction DeleteOne de la collection "event" pour supprimer l'événement correspondant à l'ID donné
	return nil
}

func GetAllComments(db *db) ([]Comment, error) {
	// Utilisez la fonction Find de la collection "comment" pour récupérer tous les commentaires dans MongoDB
	// et décodez les résultats dans une slice de Comment
	return nil, nil
}

func DeleteComment(db *db, comment_id string) error {
	// Utilisez la fonction DeleteOne de la collection "comment" pour supprimer le commentaire correspondant à l'ID donné
	return nil
}
