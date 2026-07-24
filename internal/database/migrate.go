package database

//create tables
import (
	"database/sql"
)

const schema = `

`

func Migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
