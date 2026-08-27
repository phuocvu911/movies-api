package database

//used for open connection to sqlite3 database
import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Open the connection to data.db and set maxopenconnection is 1
func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data.db?_foreign_keys=on") //turn on foreign keys in schema
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1) //1 reader, 1 writer at a time, might change it later to concurency readers after everthing done

	return db, nil
}
