package cookies

import (
	"hackaton/app/config"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Verify if a cookie exist, add a boolean Value to the map
func GetCookie(w http.ResponseWriter, data_info map[string]interface{}, r *http.Request) {
	// get database connection and error
	db, err := config.GetDB()
	if err != nil {
		panic(err) // Gérer l'erreur selon vos besoins
	}
	defer db.CloseDB() // Assurez-vous que la connexion se ferme lorsque la fonction main() se termine

	// get cookie
	cookie, err := r.Cookie("session")
	if err == nil && data_info["user"] != "" {
		data_info["cookieExist"] = true
		data_info["username"] = data_info["user"]
		data_info["connectedUserId"] = config.GetUserID(db, data_info["user"].(string))
		// Recréé un cookie pour réinitialiser le temps d'expéritation
		id := uuid.NewV4()
		cookie = &http.Cookie{
			Name:     "session",   // nom du cookie
			Value:    id.String(), // uuid pour le cookie
			HttpOnly: true,        // protection pour que le cookie ne soit pas visible par le JS
			Path:     "/",         // cookie valable de puis la racine du serveur
			MaxAge:   60 * 5,      // cookie valable 5 minutes
		}
		http.SetCookie(w, cookie)
	} else {
		// Vide tous les données concernant le précédent utilisateur connecté
		data_info["cookieExist"] = false
		data_info["username"] = ""
		data_info["role"] = ""
		data_info["user"] = ""
	}
}
