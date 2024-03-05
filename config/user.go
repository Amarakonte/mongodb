package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func AddUser(input_username string, input_email string, input_password string, info map[string]interface{}) {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}

	// Vérifiez si l'utilisateur existe déjà dans la base de données
	var existingUser User
	err = db.Database.Collection("user").FindOne(context.Background(), bson.M{"username": input_username}).Decode(&existingUser)
	if err == nil {
		info["credentials_used"] = true
		return
	}

	err = db.Database.Collection("user").FindOne(context.Background(), bson.M{"email": input_email}).Decode(&existingUser)
	if err == nil {
		info["credentials_used"] = true
		return
	}

	// Déterminez le rôle de l'utilisateur en fonction du nombre d'utilisateurs dans la base de données
	var roleID int
	count, err := db.Database.Collection("user").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		roleID = 2
	} else {
		roleID = 1
	}

	// Insérez le nouvel utilisateur dans la base de données
	user := User{
		Username: input_username,
		Email:    input_email,
		Password: input_password,
		RoleID:   roleID,
	}

	_, err = db.Database.Collection("user").InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
	}

	info["accountCreated"] = true
}

// Les autres fonctions de users.go restent les mêmes
func GetUserID(db *db, username string) string {
	var user User
	err := db.Database.Collection("user").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		// Gérer l'erreur, par exemple en renvoyant une chaîne vide
		log.Println("Erreur lors de la recherche de l'utilisateur:", err)
		return ""
	}

	return user.Id
}

func GetUser(db *db, data map[string]interface{}, user_id string) User {
	var user User
	err := db.Database.Collection("user").FindOne(context.Background(), bson.M{"_id": user_id}).Decode(&user)
	if err != nil {
		log.Println("Erreur lors de la recherche de l'utilisateur:", err)
		// Gérer l'erreur, par exemple en renvoyant un utilisateur vide
		return User{}
	}

	return user
}

func GetAllUsers(db *db) []User {
	var users []User
	cursor, err := db.Database.Collection("user").Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("Erreur lors de la recherche de tous les utilisateurs:", err)
		// Gérer l'erreur, par exemple en renvoyant une liste d'utilisateurs vide
		return []User{}
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Println("Erreur lors du décodage de l'utilisateur:", err)
			// Gérer l'erreur, par exemple en continuant à parcourir les résultats
			continue
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Erreur lors du parcours des résultats des utilisateurs:", err)
		// Gérer l'erreur, par exemple en renvoyant une liste d'utilisateurs incomplète
		return users
	}

	return users
}

// Les autres fonctions restent similaires, en utilisant les méthodes appropriées de mongo.Database

func DeleteUser(db *db, user_id string) {
	_, err := db.Database.Collection("user").DeleteOne(context.Background(), bson.M{"_id": user_id})
	if err != nil {
		panic(err)
	}
}

func MakeAdmin(db *db, userID string) {
	_, err := db.Database.Collection("user").UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"roleID": 2}})
	if err != nil {
		panic(err)
	}
}
