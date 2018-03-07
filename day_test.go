package orthocal_test

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"orthocal"
	"testing"
)

func TestDB(t *testing.T) {
	db, e := sql.Open("sqlite3", "oca_calendar.db")
	if e != nil {
		t.Errorf("Got error opening database: %#n.", e)
	}

	bibledb, e := sql.Open("sqlite3", "kjv.db")
	if e != nil {
		t.Errorf("Got error opening database: %#n.", e)
	}
	bible := orthocal.NewBible(bibledb)

	factory := orthocal.NewDayFactory(false, true, db)

	// Sunday of the Publican and Pharisee
	// Reserves should be: 266, 161, 168
	// ExtraSundays should be 3
	day := factory.NewDay(2018, 1, 28, bible)
	actual, _ := json.MarshalIndent(day, "", "\t")
	t.Errorf("%s", actual)

	// Cheesefare Sunday
	day = factory.NewDay(2018, 2, 18, bible)
	actual, _ = json.MarshalIndent(day, "", "\t")
	t.Errorf("%s", actual)

	// Annunciation
	day = factory.NewDay(2018, 3, 25, nil)
	actual, _ = json.MarshalIndent(day, "", "\t")
	t.Errorf("%s", actual)

	t.Fail()
}
