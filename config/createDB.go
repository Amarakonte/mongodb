package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CreateDB() (*sql.DB, error) {
	// Créez la chaîne de connexion à la base de données
	dsn := fmt.Sprintf("%s@tcp(%s:%s)/", userDB, ip, port)

	// Ouvrez une connexion à la base de données
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// statementa, err := db.Prepare("DROP DATABASE IF EXISTS hackaton")
	// if err != nil {
	// 	panic(err)
	// }
	// statementa.Exec()

	statement, err := db.Prepare("CREATE DATABASE IF NOT EXISTS hackaton")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	// Fermez la connexion à la base de données principale
	db.Close()

	// Ouvrez une nouvelle connexion à la base de données 'hackaton'
	db, err = sql.Open("mysql", fmt.Sprintf("%s@tcp(%s:%s)/hackaton", userDB, ip, port))
	if err != nil {
		return nil, err
	}

	return db, nil

}

func (databases db) CreateRoleTable() {
	statement, err := databases.Database.Prepare("CREATE TABLE IF NOT EXISTS role (id INTEGER PRIMARY KEY AUTO_INCREMENT, name TEXT)")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	// range over the database and check if there is double username/email
	rows, err := databases.Database.Query("SELECT COUNT(*) FROM role")
	if err != nil {
		panic(err)
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
		//stop the function if a double is found
	}
	rows.Close()

	if count < 1 {
		// add the inputs to the database
		tx, err := databases.Database.Begin()
		if err != nil {
			panic(err)
		}
		stmt, err := tx.Prepare("INSERT INTO role (name) VALUES (?)")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec("USER")
		if err != nil {
			panic(err)
		}
		tx.Commit()

		// add the inputs to the database
		tx, err = databases.Database.Begin()
		if err != nil {
			panic(err)
		}
		stmt, err = tx.Prepare("INSERT INTO role (name) VALUES (?)")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec("ADMIN")
		if err != nil {
			panic(err)
		}
		tx.Commit()
	}

}
func (databases db) CreateUserTable() {
	statement, err := databases.Database.Prepare(`
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			username TEXT,
			email TEXT,
			password TEXT,
			roleID INTEGER NOT NULL,
			FOREIGN KEY (roleID) REFERENCES role(id)
		)
	`)
	if err != nil {
		panic(err)
	}
	statement.Exec()

	// Vérifiez si l'utilisateur existe déjà
	var count int
	err = databases.Database.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", "Paris Event").Scan(&count)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		// L'utilisateur n'existe pas, alors insérez-le
		_, err = databases.Database.Exec(`
			INSERT INTO user (username, email, password, roleID)
			VALUES (?, ?, ?, ?)
		`, "Paris Event", "www.paris.fr", "test", 2)
		if err != nil {
			panic(err)
		}
	}
}

func (databases db) CreateEventTable() {
	statement, err := databases.Database.Prepare("CREATE TABLE IF NOT EXISTS event (id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255), description TEXT, date_start DATETIME, date_end DATETIME, url VARCHAR(255), cover_url VARCHAR(255) NULL, creatorID INT, FOREIGN KEY (creatorID) REFERENCES user(id) ON DELETE CASCADE, note numeric(3,2) NOT NULL, nbVote INTEGER NOT NULL)")
	if err != nil {
		panic(err)
	}
	statement.Exec()
}

func (databases db) CreateParticipentsTable() {
	statement, err := databases.Database.Prepare("CREATE TABLE IF NOT EXISTS participants (eventID INTEGER NOT NULL, FOREIGN KEY (eventID) REFERENCES event(id) ON DELETE CASCADE, userID INTEGER NOT NULL, FOREIGN KEY (userID) REFERENCES user(id) ON DELETE CASCADE, accepted BOOL)")
	if err != nil {
		panic(err)
	}
	statement.Exec()
}

func (databases db) CreateCommentTable() {
	statement, err := databases.Database.Prepare("CREATE TABLE IF NOT EXISTS comment (id INTEGER PRIMARY KEY AUTO_INCREMENT, content TEXT, creationDate DATETIME, userID INTEGER NOT NULL, FOREIGN KEY (userID) REFERENCES user(id) ON DELETE CASCADE, eventID INTEGER NOT NULL, FOREIGN KEY (eventID) REFERENCES event(id) ON DELETE CASCADE)")
	if err != nil {
		panic(err)
	}
	statement.Exec()
}
