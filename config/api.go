package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Api() {
	// Initialisez la connexion à votre base de données SQL ici
	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(%s:%s)/hackaton", userDB, ip, port))
	if err != nil {
		fmt.Println("Erreur lors de la connexion à la base de données:", err)
		return
	}
	defer db.Close()

	// Récupéré les données
	url := "https://parisdata.opendatasoft.com/api/explore/v2.1/catalog/datasets/que-faire-a-paris-/records?select=id%2C%20url%2C%20title%2C%20description%2C%20date_start%2C%20date_end%2C%20cover_url"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer resp.Body.Close()

	// Décodage des données
	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Erreur lors de la décodage JSON:", err)
		return
	}

	// Insérer les données dans la bdd
	for _, event := range response.Results {
		var creatorID int = 1
		_, err := db.Exec("INSERT IGNORE INTO event (id, title, description, date_start, date_end, url, cover_url, creatorID) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			event.Id, event.Title, event.Description, event.Date_start, event.Date_end, event.Url, event.Cover_url, creatorID)
		if err != nil {
			fmt.Println("Erreur lors de l'insertion dans la base de données:", err)
		}
	}

	// Fermez la connexion à la base de données lorsque vous avez terminé avec elle.
	db.Close()
}
