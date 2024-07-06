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

func GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := db.FetchItems()

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusOK, items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'id' must be of type int" })
		return
	}

	item, err := db.FetchItem(id)

	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithJSON(w, http.StatusUnprocessableEntity, ResponseJson{ "message": fmt.Sprintf("item with id %d does not exist!", id) })
		} else {
			RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var requestBody db.CreateItemBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": err.Error() })
		return
	}
	defer r.Body.Close()

	if requestBody.Item == "" {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'item' field is required!" })
		return
	}

	item, err := db.InsertItem(requestBody.Item)
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusCreated, item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var requestBody db.UpdateItemBody

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

	item, err := db.UpdateItem(id, requestBody)

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusOK, item)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, ResponseJson{ "message": "'id' must be of type int" })
		return
	}

	err = db.DeleteItem(id)

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func DeleteCompletedItems(w http.ResponseWriter, r *http.Request) {
	err := db.DeleteCompletedItems()

	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ResponseJson{ "message": err.Error() })
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
