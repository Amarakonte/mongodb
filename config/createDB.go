package config

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateDB() (*mongo.Database, error) {
	// Connexion à MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Vérification de la connexion
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Sélection de la base de données 'hackaton'
	database := client.Database("hackaton")

	return database, nil
}

func (databases *db) CreateRoleTable() error {
	// Insertion des rôles dans la collection 'role'
	roles := []interface{}{
		bson.M{"name": "USER"},
		bson.M{"name": "ADMIN"},
	}

	_, err := databases.Database.Collection("role").InsertMany(context.Background(), roles)
	if err != nil {
		return err
	}

	return nil
}

func (databases *db) CreateUserTable() error {
	// Insérer les données initiales dans la collection 'user'
	users := []interface{}{
		bson.M{"username": "Paris Event", "email": "www.paris.fr", "password": "test", "roleID": 2},
	}

	_, err := databases.Database.Collection("user").InsertMany(context.Background(), users)
	if err != nil {
		return err
	}

	return nil
}

func (databases *db) CreateEventTable() error {
	// Pas besoin de créer une table dans MongoDB
	return nil
}

func (databases *db) CreateParticipantsTable() error {
	// Pas besoin de créer une table dans MongoDB
	return nil
}

func (databases *db) CreateCommentTable() error {
	// Pas besoin de créer une table dans MongoDB
	return nil
}
