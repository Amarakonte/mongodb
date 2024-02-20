package config

func AddEvent(input_title string, input_description string, input_date string, data map[string]interface{}) {
	db := GetDB()

	query := "INSERT INTO event (title, description, date, creatorID) VALUES (?, ?, ?, ?)"

	tx, err := db.Database.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(input_title, input_description, input_date, GetUserID(db, data["username"].(string)))
	if err != nil {
		panic(err)
	}
	tx.Commit()
	data["eventCreated"] = true
}
