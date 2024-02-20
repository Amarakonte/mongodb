package config

import (
	"database/sql"
	"strings"
)

func DisplayEvents(data map[string]interface{}, searching string) {
	db := GetDB()

	var rows *sql.Rows
	var err error
	if searching == "" {
		rows, err = db.Database.Query("SELECT event.id, event.title, event.description, event.date_start, event.date_end, event.cover_url, event.url, user.username, user.id, event.note, event.nbVote FROM event INNER JOIN user ON user.id = event.creatorID")
	} else {
		rows, err = db.Database.Query("SELECT event.id, event.title, event.description, event.date_start, event.date_end, event.cover_url, event.url, user.username, user.id, event.note, event.nbVote FROM event INNER JOIN user ON user.id = event.creatorID WHERE event.title = ?", searching)
	}
	if err != nil {
		panic(err)
	}

	var events []Event

	for rows.Next() {
		var event Event
		rows.Scan(&event.Id, &event.Title, &event.Description, &event.Date_start, &event.Date_end, &event.Cover_url, &event.Url, &event.User.Username, &event.User.Id, &event.Note, &event.NbVote)

		// Remplace les \n par des <br> pour sauter des lignes en html
		event.Description = strings.Replace(strings.Replace(event.Description, string('\r'), "", -1), string('\n'), "<br>", -1)

		events = append(events, event)
	}

	defer rows.Close()

	data["events"] = events
}

func GetEvent(data map[string]interface{}, id_event string) {
	db := GetDB()

	rows, err := db.Database.Query("SELECT event.title, event.description, event.date_start, event.date_end, event.cover_url, event.url, user.username, user.id, event.note, event.nbVote FROM event INNER JOIN user ON user.id = event.creatorID WHERE event.id = ?", id_event)
	if err != nil {
		panic(err)
	}
	var event Event
	for rows.Next() {
		err := rows.Scan(&event.Title, &event.Description, &event.Date_start, &event.Date_end, &event.Cover_url, &event.Url, &event.User.Username, &event.User.Id, &event.Note, &event.NbVote)
		if err != nil {
			panic(err)
		}
	}

	if event.Title != "" {
		data["event"] = event
	}

	defer rows.Close()

	data["Id_event"] = id_event

}

func GetComments(db db, id_event string, data map[string]interface{}) {
	rows, err := db.Database.Query("SELECT comment.content, comment.creationDate, user.username, user.id FROM comment INNER JOIN user ON user.id = comment.userID  WHERE eventID = ?", id_event)
	if err != nil {
		panic(err)
	}

	var comments []Comment

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.Content, &comment.CreationDate, &comment.User.Username, &comment.User.Id)

		// Remplace les \n par des <br> pour sauter des lignes en html
		comment.Content = strings.Replace(strings.Replace(comment.Content, string('\r'), "", -1), string('\n'), "<br>", -1)

		comments = append(comments, comment)
	}

	defer rows.Close()

	data["comments"] = comments
}
