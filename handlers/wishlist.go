package handlers

import (
	"net/http"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/views/wishlist"
)

func HandleGetWishes(w http.ResponseWriter, r *http.Request) error {
	var data []db.RecordsDbRow

	rows, err := db.Db.Query("SELECT * FROM records ORDER BY createdat DESC")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var record db.RecordsDbRow
		err := rows.Scan(&record.Id, &record.Text, &record.CreatedAt, &record.UpdatedAt, &record.Completed)
		if err != nil {
			return err
		}

		data = append(data, record)
	}

	return wishlist.Wishlist(data).Render(r.Context(), w)
}

func HandleCreateWish(w http.ResponseWriter, r *http.Request) error {
	// text := r.FormValue("text")
	var data []db.RecordsDbRow
	text := "hehe"

	_, err := db.Db.Exec("INSERT INTO records (text) VALUES ($1)", text)
	if err != nil {
		return err
	}

	rows, err := db.Db.Query("SELECT * FROM records ORDER BY createdat DESC")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var record db.RecordsDbRow
		err := rows.Scan(&record.Id, &record.Text, &record.CreatedAt, &record.UpdatedAt, &record.Completed)
		if err != nil {
			return err
		}

		data = append(data, record)
	}

	return wishlist.Wishlist(data).Render(r.Context(), w)
}