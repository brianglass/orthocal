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

	// Veneration of the Cross - should include 7th Matins Gospel
	day = factory.NewDay(2018, 3, 11, nil)
	actual, _ = json.MarshalIndent(day, "", "\t")
	t.Errorf("%s", actual)

	// Memorial Saturday with no memorial readings
	// Should not have John 5.24-30
	day = factory.NewDay(2022, 3, 26, nil)
	actual, _ = json.MarshalIndent(day, "", "\t")
	t.Errorf("%s", actual)

	/*
		today := time.Now()
		for {
			today = today.AddDate(0, 0, 1)
			day = factory.NewDay(today.Year(), int(today.Month()), today.Day(), nil)
			if day.HasNoMemorial() {
				actual, _ = json.MarshalIndent(day, "", "\t")
				t.Errorf("%s", actual)
				break
			}
		}
	*/

	t.Fail()
}
