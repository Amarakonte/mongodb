package cookies

import (
	"context"
	"hackaton/app/config"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func SearchUserToLog(data_logIn map[string]interface{}, user_login string, password_login string, data map[string]interface{}) bool {
	var create_cookie = false

	db, err := config.GetDB()
	if err != nil {
		panic(err)
	}
	defer db.CloseDB() // Assurez-vous que la connexion se ferme lorsque la fonction main() se termine

	// Parcourir la BDD
	rows, err := db.Database.Collection("user").Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	defer rows.Close(context.Background())
	for rows.Next(context.Background()) {
		var user config.User
		err := rows.Decode(&user)
		if err != nil {
			panic(err)
		}

		if user_login == user.Username && password_login == user.Password {
			create_cookie = true
			data["user"] = user.Username
			data["role"] = user.Role.Name
			data["cookieExist"] = true
			break
		} else if user_login != "" {
			data_logIn["wrongCredentials"] = true
		}
	}
	return create_cookie
}

func SetDataToSend(w http.ResponseWriter, r *http.Request, data_info map[string]interface{}, data map[string]interface{}, on_user_page bool, user_page string) {
	// Copiez la carte principale pour obtenir toutes les informations importantes
	for k, v := range data {
		data_info[k] = v
	}
	data_info["cookieExist"] = false
	data_info["username"] = ""
	GetCookie(w, data_info, r) // ./cookies/getCookies.go
}
