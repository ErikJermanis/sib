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

	return record, err
}

func FetchOtpDetails(otp string) (OtpsDbRow, error) {
	var otpDetails OtpsDbRow

	err := Db.QueryRow("SELECT used, expiresat FROM otps WHERE otp = $1", otp).Scan(&otpDetails.Used, &otpDetails.ExpiresAt)

	return otpDetails, err
}

func InvalidateOtp(otp string) error {
	_, err := Db.Exec("UPDATE otps SET used = true WHERE otp = $1", otp)
	return err
}

func InsertRecord(text string) (RecordsDbRow, error) {
	var wish RecordsDbRow

	err := Db.QueryRow("INSERT INTO records (text) VALUES ($1) RETURNING *", text).Scan(&wish.Id, &wish.Text, &wish.CreatedAt, &wish.UpdatedAt, &wish.Completed)

	return wish, err
}

func UpdateRecord(id int, body UpdateRecordBody) (RecordsDbRow, error) {
	var wish RecordsDbRow
	var err error

	if (body.Text != "") {
		err = Db.QueryRow("UPDATE records SET text = $1, completed = $2, updatedat = NOW() WHERE id = $3 RETURNING *", body.Text, body.Completed, id).Scan(&wish.Id, &wish.Text, &wish.CreatedAt, &wish.UpdatedAt, &wish.Completed)
	} else {
		err = Db.QueryRow("UPDATE records SET completed = $1, updatedat = NOW() WHERE id = $2 RETURNING *", body.Completed, id).Scan(&wish.Id, &wish.Text, &wish.CreatedAt, &wish.UpdatedAt, &wish.Completed)
	}

	return wish, err
}

func DeleteRecord(id int) error {
	_, err := Db.Exec("DELETE FROM records WHERE id = $1", id)
	return err
}

func FetchItems() ([]ItemsDbRow, error) {
	var items []ItemsDbRow

	rows, err := Db.Query("SELECT * FROM items ORDER BY (CASE WHEN completed THEN 0 ELSE 1 END) DESC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row ItemsDbRow
		if err = rows.Scan(&row.Id, &row.Item, &row.Completed, &row.Rank); err != nil {
			return nil, err
		}
		items = append(items, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func FetchItem(id int) (ItemsDbRow, error) {
	var item ItemsDbRow

	row := Db.QueryRow("SELECT * FROM items WHERE id = $1", id)
	
	err := row.Scan(&item.Id, &item.Item, &item.Completed, &item.Rank)

	return item, err
}

func InsertItem(itemContent string) (ItemsDbRow, error) {
	var item ItemsDbRow
	var count int

	row := Db.QueryRow("SELECT COUNT(*) FROM items")
	err := row.Scan(&count)

	if err != nil {
		return item, err
	}

	err = Db.QueryRow("INSERT INTO items (item, rank) VALUES ($1, $2) RETURNING *", itemContent, count + 1).Scan(&item.Id, &item.Item, &item.Completed, &item.Rank)

	return item, err
}

func UpdateItem(id int, body UpdateItemBody) (ItemsDbRow, error) {
	var item ItemsDbRow
	var err error

	if body.Item == "" && body.Rank == 0 {
		err = Db.QueryRow("UPDATE items SET completed = $1 WHERE id = $2 RETURNING *", body.Completed, id).Scan(&item.Id, &item.Item, &item.Completed, &item.Rank)
	} else if body.Item == "" {
		err = Db.QueryRow("UPDATE items SET completed = $1, rank = $2 WHERE id = $3 RETURNING *", body.Completed, body.Rank, id).Scan(&item.Id, &item.Item, &item.Completed, &item.Rank)
	} else if body.Rank == 0 {
		err = Db.QueryRow("UPDATE items SET item = $1, completed = $2 WHERE id = $3 RETURNING *", body.Item, body.Completed, id).Scan(&item.Id, &item.Item, &item.Completed, &item.Rank)
	} else {
		err = Db.QueryRow("UPDATE items SET item = $1, completed = $2, rank = $3 WHERE id = $4 RETURNING *", body.Item, body.Completed, body.Rank, id).Scan(&item.Id, &item.Item, &item.Completed, &item.Rank)
	}

	return item, err
}

func DeleteItem(id int) error {
	_, err := Db.Exec("DELETE FROM items WHERE id = $1", id)
	return err
}

func DeleteCompletedItems() error {
	_, err := Db.Exec("DELETE FROM items WHERE completed = true")
	return err
}
