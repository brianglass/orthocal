package orthocal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	Driver   = "sqlite3"
	Database = "orthocal.sqlite"
)

func TestDB() {
	db, e := sql.Open(Driver, Database)
	if e != nil {
		log.Printf("Got error opening database: %#n.", e)
	}

	rows, e := db.Query("select pdist, title, subtitle from days")
	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
	}

	for rows.Next() {
		var pdist int16
		var title, subtitle string

		rows.Scan(&pdist, &title, &subtitle)
		if len(subtitle) > 0 {
			fmt.Printf("pdist = %d, title = \"%s: %s\"\n", pdist, title, subtitle)
		} else {
			fmt.Printf("pdist = %d, title = \"%s\"\n", pdist, title)
		}
	}
}
