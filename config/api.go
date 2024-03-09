package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Api() {
	ip := "localhost"
	port := "27017"

	// Connexion à MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://" + ip + ":" + port)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Erreur lors de la connexion à la base de données MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Récupération des données depuis l'API
	url := "https://parisdata.opendatasoft.com/api/explore/v2.1/catalog/datasets/que-faire-a-paris-/records?select=id%2C%20url%2C%20title%2C%20description%2C%20date_start%2C%20date_end%2C%20cover_url"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer resp.Body.Close()

	// Décodage des données JSON
	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	// Insérer les données dans la collection "event" de la base de données MongoDB
	eventCollection := client.Database("hackaton").Collection("event")
	for _, event := range response.Results {
		creatorID := "2" // Définir l'ID du créateur de l'événement
		event.CreatorID = creatorID
		event.Timestamp = time.Now().Unix()

		// Utiliser l'opérateur upsert pour mettre à jour ou insérer le document
		filter := bson.M{"id": event.Id} // Utilisez un champ unique comme filtre (id dans ce cas)
		update := bson.M{"$set": event}  // Mettre à jour avec les données de l'événement
		_, err := eventCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
		if err != nil {
			fmt.Println("Erreur lors de l'insertion/mise à jour dans la base de données MongoDB:", err)
		}
	}

}
