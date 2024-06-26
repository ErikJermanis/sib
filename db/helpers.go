package db

func FetchRecords() ([]RecordsDbRow, error) {
	var records []RecordsDbRow

	rows, err := Db.Query("SELECT * FROM records ORDER BY createdat DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row RecordsDbRow
		if err = rows.Scan(&row.Id, &row.Text, &row.CreatedAt, &row.UpdatedAt, &row.Completed); err != nil {
			return nil, err
		}
		records = append(records, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func FetchRecord(id int) (RecordsDbRow, error) {
	var record RecordsDbRow

	row := Db.QueryRow("SELECT * FROM records WHERE id = $1", id)
	
	err := row.Scan(&record.Id, &record.Text, &record.CreatedAt, &record.UpdatedAt, &record.Completed)
	if err != nil {
		return record, err
	}

	return record, nil
}
