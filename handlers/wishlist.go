package handlers

import (
	"fmt"
	"net/http"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/views/wishlist"
)

func HandleGetWishes(w http.ResponseWriter, r *http.Request) error {
	rows, err := db.Db.Query("SELECT * FROM records ORDER BY createdat DESC")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var record db.RecordsRow
		err := rows.Scan(&record.Id, &record.Text, &record.CreatedAt, &record.UpdatedAt, &record.Completed)
		if err != nil {
			return err
		}

		fmt.Println(record)
	}

	return wishlist.Wishlist().Render(r.Context(), w)
}