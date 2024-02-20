package config

import (
	"log"
)

var username, email string

func AddUser(input_username string, input_email string, input_password string, info map[string]interface{}) {
	db := GetDB()

	// range over the database and check if there is double username/email
	rows, err := db.Database.Query("SELECT username, email FROM user")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&username, &email)
		//stop the function if a double is found
		if username == input_username {
			info["credentials_used"] = true
			rows.Close()
			return
		} else if email == input_email {
			info["credentials_used"] = true
			rows.Close()
			return
		}
	}

	rows.Close()

	// Compte le nb d'utiliteur, le 1er sera Admin
	rows, err = db.Database.Query("SELECT COUNT(*) FROM user")
	if err != nil {
		log.Fatal(err)
	}
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	defer rows.Close()

	query := "INSERT INTO user (username, email, password, roleID) VALUES (?, ?, ?, ?)"
	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		_, err = stmt.Exec(input_username, input_email, input_password, 1)
	} else {
		_, err = stmt.Exec(input_username, input_email, input_password, 2)
	}
	if err != nil {
		panic(err)
	}
	tx.Commit()
	info["accountCreated"] = true

	// // Test d'envoie de mail
	// auth := smtp.PlainAuth("", "user@example.com", "password", "mail.example.com")

	// to := []string{input_email}
	// msg := []byte("To: " + input_email + "\r\n" +
	// 	"Sujet: Inscription terminée!\r\n" +
	// 	"\r\n" +
	// 	"This is the email body.\r\n")
	// err = smtp.SendMail("mail.example.com:1025", auth, "sender@example.org", to, msg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func GetUserID(db db, username string) string {
	rows, err := db.Database.Query("SELECT id FROM user WHERE username = ?", username)
	if err != nil {
		panic(err)
	}
	var user_id string
	for rows.Next() {
		rows.Scan(&user_id)
	}
	rows.Close()
	return user_id
}

func GetUser(db db, data map[string]interface{}, user_id string) User {
	rows, err := db.Database.Query("SELECT user.username, user.email, role.name, SUM(event.note)/COUNT(event.id) FROM user INNER JOIN event ON event.creatorID = user.id INNER JOIN role ON role.id = user.roleID WHERE user.id = ?", user_id)
	if err != nil {
		panic(err)
	}
	var user User
	for rows.Next() {
		rows.Scan(&user.Username, &user.Email, &user.Role.Name, &user.Moyenne)
	}
	rows.Close()

	return user
}

func GetAllUsers(db db) []User {
	rows, err := db.Database.Query("SELECT user.id, user.username, user.email, role.name FROM `user` INNER JOIN role ON role.id = user.roleID ORDER BY user.id")
	if err != nil {
		panic(err)
	}
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role.Name)
		users = append(users, user)
	}
	rows.Close()

	return users
}

func DeleteUser(db db, user_id string) {
	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	// Supprimer la ligne de l'utilisateur dans la BDD users
	stmt, err := tx.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(user_id)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func MakeAdmin(db db, userID string) {
	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	query := "UPDATE user SET roleID = 2 WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(userID)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
