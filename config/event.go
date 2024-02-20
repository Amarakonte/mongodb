package config

import (
	"strconv"
	"time"
)

func AddCommentOnEvent(id_event string, input_commentaire string, data_PageEvent map[string]interface{}) {
	db := GetDB()

	query := "INSERT INTO comment (content, creationDate, userID, eventID) VALUES (?, ?, ?, ?)"

	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(input_commentaire, time.Now(), GetUserID(db, data_PageEvent["username"].(string)), id_event)
	if err != nil {
		panic(err)
	}
	tx.Commit()
	data_PageEvent["commentaireCreated"] = true
}

func UpdateNote(db db, data map[string]interface{}, id_event string, newNote string) {
	rows, err := db.Database.Query("SELECT note, nbVote FROM event WHERE id = ?", id_event)
	if err != nil {
		panic(err)
	}

	var noteString, nbVoteString string
	for rows.Next() {
		rows.Scan(&noteString, &nbVoteString)
	}

	noteFloat, err := strconv.ParseFloat(noteString, 32)
	if err != nil {
		panic(err)
	}

	newNoteFloat, err := strconv.ParseFloat(newNote, 32)
	if err != nil {
		panic(err)
	}

	nbVoteInt, err := strconv.Atoi(nbVoteString)
	if err != nil {
		panic(err)
	}

	noteFinal := ((noteFloat * float64(nbVoteInt)) + newNoteFloat) / (float64(nbVoteInt) + 1)

	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	query := "UPDATE event SET note = ?, nbVote = ? WHERE id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(noteFinal, nbVoteInt+1, id_event)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	tx.Commit()
}

func GetAllEvents(db db) []Event {
	rows, err := db.Database.Query("SELECT event.id, event.title, event.description, event.date_start, event.date_end, event.cover_url, user.username, user.id FROM event INNER JOIN user ON user.id = event.creatorID ORDER BY event.id")
	if err != nil {
		panic(err)
	}
	var events []Event
	for rows.Next() {
		var event Event
		rows.Scan(&event.Id, &event.Title, &event.Description, &event.Date_start, &event.Date_end, &event.Cover_url, &event.User.Username, &event.User.Id)
		events = append(events, event)
	}

	rows.Close()

	return events
}

func DeleteEvent(db db, event_id string) {
	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	// Supprimer la ligne de l'utilisateur dans la BDD users
	stmt, err := tx.Prepare("DELETE FROM event WHERE id = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(event_id)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func GetAllComments(db db) []Comment {
	rows, err := db.Database.Query("SELECT comment.id, comment.content, comment.creationDate, event.title, user.username FROM comment INNER JOIN user ON user.id = comment.userID INNER JOIN event ON event.id = comment.eventID ORDER BY event.id")
	if err != nil {
		panic(err)
	}
	var comments []Comment
	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.Id, &comment.Content, &comment.CreationDate, &comment.Event.Title, &comment.User.Username)
		comments = append(comments, comment)
	}
	rows.Close()

	return comments
}

func DeleteComment(db db, comment_id string) {
	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	// Supprimer la ligne de l'utilisateur dans la BDD users
	stmt, err := tx.Prepare("DELETE FROM comment WHERE id = ?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(comment_id)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
