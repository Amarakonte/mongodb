package config

func AddParticipant(id_event string, data map[string]interface{}) {
	db := GetDB()

	query := "INSERT INTO participants (eventID, userID, accepted) VALUES (?, ?, ?)"

	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(id_event, GetUserID(db, data["username"].(string)), false)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func GetParticipants(db db, id_event string, data map[string]interface{}) {
	rows, err := db.Database.Query("SELECT user.username, participants.accepted, participants.userID, participants.eventID FROM participants INNER JOIN user ON participants.userID = user.id WHERE participants.eventID = ?", id_event)
	if err != nil {
		panic(err)
	}

	var participants []Participants

	for rows.Next() {
		var participant Participants
		rows.Scan(&participant.User.Username, &participant.Accepted, &participant.User.Id, &participant.Event.Id)

		participants = append(participants, participant)
	}

	defer rows.Close()

	data["participants"] = participants

}

func AcceptParticipation(db db, id_event string, id_user string) {

	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	query := "UPDATE participants SET accepted = 1 WHERE eventID = ? AND userID = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(id_event, id_user)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func RemoveParticipant(db db, id_event string, id_user string) {
	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	query := "DELETE FROM participants WHERE eventID = ? AND userID = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(id_event, id_user)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func IsParticipant(db db, id_event string, data map[string]interface{}) bool {

	rows, err := db.Database.Query("SELECT accepted FROM participants WHERE eventID = ? AND userID = ?", id_event, GetUserID(db, data["username"].(string)))
	if err != nil {
		panic(err)
	}

	var isAccepted bool
	for rows.Next() {
		rows.Scan(&isAccepted)
	}
	defer rows.Close()

	return isAccepted
}

func HasRequestedParticipation(db db, id_event string, data map[string]interface{}) bool {
	rows, err := db.Database.Query("SELECT COUNT(*) FROM participants WHERE eventID = ? AND userID = ?", id_event, GetUserID(db, data["username"].(string)))
	if err != nil {
		panic(err)
	}

	var requested int
	for rows.Next() {
		rows.Scan(&requested)
	}

	defer rows.Close()

	return requested > 0
}
