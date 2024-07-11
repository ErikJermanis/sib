package handlers

import (
	"net/http"
	"strconv"

	"github.com/ErikJermanis/sib-web/db"
	"github.com/ErikJermanis/sib-web/views/shoplist"
	"github.com/go-chi/chi/v5"
)

func HandleGetItems(w http.ResponseWriter, r *http.Request) error {
	data, err := db.FetchItems()

	if err != nil {
		return err
	}

	var containsCompleted bool = false
	if len(data) > 0 {
		containsCompleted = data[len(data)-1].Completed
	}

	return shoplist.Index(data, containsCompleted).Render(r.Context(), w)
}

func HandleCreateItem(w http.ResponseWriter, r *http.Request) error {
	item := r.FormValue("item")

	// TODO: form validation
	// TODO: look into sending just one item and somehow retaining sorting on client

	_, err := db.InsertItem(item)
	if err != nil {
		return err
	}

	items, err := db.FetchItems()
	if err != nil {
		return err
	}

	return shoplist.NewItem(items).Render(r.Context(), w)
}

func HandleCompleteItem(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	_, err = db.UpdateItem(id, db.UpdateItemBody{ Completed: true })
	if err != nil {
		return err
	}

	items, err := db.FetchItems()
	if err != nil {
		return err
	}

	return shoplist.ModifyItem(items, true).Render(r.Context(), w)
}

func HandleResetItem(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	_, err = db.UpdateItem(id, db.UpdateItemBody{ Completed: false })
	if err != nil {
		return err
	}

	items, err := db.FetchItems()
	if err != nil {
		return err
	}
	var containsCompleted bool = false
	if len(items) > 0 {
		containsCompleted = items[len(items)-1].Completed
	}

	return shoplist.ModifyItem(items, containsCompleted).Render(r.Context(), w)
}

func HandleDeleteCompletedItems(w http.ResponseWriter, r *http.Request) error {
	err := db.DeleteCompletedItems()
	if err != nil {
		return err
	}

	items, err := db.FetchItems()
	if err != nil {
		return err
	}

	return shoplist.ModifyItem(items, false).Render(r.Context(), w)
}
