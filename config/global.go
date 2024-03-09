package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ip   = "localhost"
	port = "27017"
)

type APIResponse struct {
	TotalCount int     `json:"total_count"`
	Results    []Event `json:"results"`
}

type Event struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date_start  string `json:"date_start"`
	Date_end    string `json:"date_end"`
	URL         string `json:"url"`
	Cover_url   string `json:"cover_url"`
	CreatorID   string
	Timestamp   int64
	User        User
	Note        string
	NbVote      string
}

type Comment struct {
	Id           string `json:"id"`
	Content      string
	CreationDate string
	User         User
	Event        Event
}

type User struct {
	Id       string
	Email    string
	Username string
	Moyenne  string
	Role     Role
	Events   []Event
	Password string
	RoleID   int
}

type Participants struct {
	Event    Event
	User     User
	Accepted bool
}

type Role struct {
	Name string
}

type db struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func GetDB() (*db, error) {
	// Connexion à MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://" + ip + ":" + port)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Vérification de la connexion
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	database := client.Database("hackaton")

	return &db{
		Client:   client,
		Database: database,
	}, nil
}

// Fermez la connexion à la base de données lorsque vous n'en avez plus besoin
func (d *db) CloseDB() {
	if d.Client != nil {
		d.Client.Disconnect(context.Background())
	}
}
