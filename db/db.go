package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type RecordsDbRow struct {
	Id int `json:"id"`
	Text string `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Completed bool `json:"completed"`
}

type CreateRecordBody struct {
	Text string
}

type UpdateRecordBody struct {
	Text string
	Completed bool
}

type OtpsDbRow struct {
	Used bool
	ExpiresAt time.Time
}

type Record struct {
	Text string
	Completed bool
}

type ItemsDbRow struct {
	Id int `json:"id"`
	Item string `json:"item"`
	Completed bool `json:"completed"`
	Rank int `json:"rank"`
}

type CreateItemBody struct {
	Item string
}

type UpdateItemBody struct {
	Item string
	Completed bool
	Rank int
}

var Db *sql.DB

func Initialize(user, password, dbname string) {
	connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	Db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
}