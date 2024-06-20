package handlers

import (
	"net/http"
	"strconv"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/views/wishlist"
	"github.com/go-chi/chi/v5"
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

	return wishlist.Index(data).Render(r.Context(), w)
}

func HandleCreateWish(w http.ResponseWriter, r *http.Request) error {
	// TODO: handle form submission
	// text := r.FormValue("text")
	var data db.RecordsDbRow
	text := "hehe"

	err := db.Db.QueryRow("INSERT INTO records (text) VALUES ($1) RETURNING *", text).Scan(&data.Id, &data.Text, &data.CreatedAt, &data.UpdatedAt, &data.Completed)
	if err != nil {
		return err
	}

	return wishlist.ListItem(data).Render(r.Context(), w)
}

func HandleSelectWish(w http.ResponseWriter, r *http.Request) error {
	var data db.RecordsDbRow

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	rows, err := db.Db.Query("SELECT * FROM records WHERE id = $1", id)
	if err != nil {
		return err
	}

	if rows.Next() {
		err := rows.Scan(&data.Id, &data.Text, &data.CreatedAt, &data.UpdatedAt, &data.Completed)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	return wishlist.ListItemSelected(data).Render(r.Context(), w)
}

func HandleDeselectWish(w http.ResponseWriter, r *http.Request) error {
	var data db.RecordsDbRow

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	rows, err := db.Db.Query("SELECT * FROM records WHERE id = $1", id)
	if err != nil {
		return err
	}

	if rows.Next() {
		err := rows.Scan(&data.Id, &data.Text, &data.CreatedAt, &data.UpdatedAt, &data.Completed)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	return wishlist.ListItem(data).Render(r.Context(), w)
}

func HandleCompleteWish(w http.ResponseWriter, r *http.Request) error {
	var data db.RecordsDbRow

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	err = db.Db.QueryRow("UPDATE records SET completed = true, updatedat = NOW() WHERE id = $1 RETURNING *", id).Scan(&data.Id, &data.Text, &data.CreatedAt, &data.UpdatedAt, &data.Completed)
	if err != nil {
		return err
	}

	return wishlist.ListItem(data).Render(r.Context(), w)
}

func HandleDeleteWish(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	_, err = db.Db.Exec("DELETE FROM records WHERE id = $1", id)
	if err != nil {
		return err
	}

	// TODO: status no content here and also improve statuses elsewhere (learn how to handle them with htmx)

	return nil
}

func HandleResetWish(w http.ResponseWriter, r *http.Request) error {
	var data db.RecordsDbRow

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	err = db.Db.QueryRow("UPDATE records SET completed = false, updatedat = NOW() WHERE id = $1 RETURNING *", id).Scan(&data.Id, &data.Text, &data.CreatedAt, &data.UpdatedAt, &data.Completed)
	if err != nil {
		return err
	}

	return wishlist.ListItem(data).Render(r.Context(), w)
}
