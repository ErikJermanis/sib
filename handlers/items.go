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

	return shoplist.Index(data).Render(r.Context(), w)
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

	return shoplist.List(items).Render(r.Context(), w)
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

	return shoplist.List(items).Render(r.Context(), w)
}
