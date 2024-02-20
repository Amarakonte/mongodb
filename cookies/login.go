package cookies

import (
	"hackaton/app/config"
	"net/http"
)

func SearchUserToLog(data_logIn map[string]interface{}, user_login string, password_login string, data map[string]interface{}) bool {
	var create_cookie = false

	db := config.GetDB()
	defer db.CloseDB() // Assurez-vous que la connexion se ferme lorsque la fonction main() se termine

	// Parcourir la BDD
	rows, err := db.Database.Query("SELECT username, password FROM user")
	if err != nil {
		panic(err)
	}

	var username, password string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &password)
		if user_login == username && password == password_login {
			user := config.GetUser(config.GetDB(), data, config.GetUserID(config.GetDB(), username))
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
	// Copy the main map to get all important info
	for k, v := range data {
		data_info[k] = v
	}
	data_info["cookieExist"] = false
	data_info["username"] = ""
	GetCookie(w, data_info, r) // ./cookies/getCookies.go
}
