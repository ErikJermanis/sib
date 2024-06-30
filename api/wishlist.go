package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/go-chi/chi/v5"
)

func GetWishes(w http.ResponseWriter, r *http.Request) {
	records, err := db.FetchRecords()

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusOK, records)
}

func GetWish(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'id' must be of type int" })
		return
	}

	record, err := db.FetchRecord(id)

	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithJSON(w, http.StatusUnprocessableEntity, ResponseJson{ "message": fmt.Sprintf("record with id %d does not exist!", id) })
		} else {
			RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, record)
}

func CreateWish(w http.ResponseWriter, r *http.Request) {
	var requestBody db.CreateRecordBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": err.Error() })
		return
	}
	defer r.Body.Close()

	if requestBody.Text == "" {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'text' field is required!" })
		return
	}

	wish, err := db.InsertRecord(requestBody.Text)
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusCreated, wish)
}

func UpdateWish(w http.ResponseWriter, r *http.Request) {
	var requestBody db.UpdateRecordBody

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'id' must be of type int" })
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": err.Error() })
		return
	}
	defer r.Body.Close()

	wish, err := db.UpdateRecord(id, requestBody)

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusOK, wish)
}

func DeleteWish(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'id' must be of type int" })
		return
	}

	err = db.DeleteRecord(id)

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
