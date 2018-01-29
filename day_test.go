package orthocal_test

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orthocal"
	"testing"
)

func TestDB(t *testing.T) {
	db, e := sql.Open("sqlite3", "oca_calendar.db")
	if e != nil {
		t.Errorf("Got error opening database: %#n.", e)
	}

	// Sunday of the Publican and Pharisee
	day := orthocal.NewDay(2018, 1, 28, false, db)
	day.GetRecords()

	// Annunciation
	day = orthocal.NewDay(2018, 3, 25, false, db)
	day.GetRecords()

	t.Fail()
}
