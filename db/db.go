package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type RecordsDbRow struct {
	Id int
	Text string
	CreatedAt time.Time
	UpdatedAt time.Time
	Completed bool
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

var Db *sql.DB

func Initialize(user, password, dbname string) {
	connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	Db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
}