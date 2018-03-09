package orthocal_test

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orthocal"
	"strings"
	"testing"
)

func TestDB(t *testing.T) {
	db, e := sql.Open("sqlite3", "oca_calendar.db")
	if e != nil {
		t.Errorf("Got error opening database: %#n.", e)
	}

	/*
		bibledb, e := sql.Open("sqlite3", "kjv.db")
		if e != nil {
			t.Errorf("Got error opening database: %#n.", e)
		}
		bible := orthocal.NewBible(bibledb)
	*/

	factory := orthocal.NewDayFactory(false, true, db)

	/*
		// Sunday of the Publican and Pharisee
		day := factory.NewDay(2018, 1, 28, bible)
		actual, _ := json.MarshalIndent(day, "", "\t")
		t.Errorf("%s", actual)

		// Cheesefare Sunday
		day = factory.NewDay(2018, 2, 18, bible)
		actual, _ = json.MarshalIndent(day, "", "\t")
		t.Errorf("%s", actual)
	*/

	t.Run("Annunciation Commemorations", func(t *testing.T) {
		day := factory.NewDay(2018, 3, 25, nil)

		count := 0
		for _, c := range day.Commemorations {
			if c.FeastName == "Annunciation Most Holy Theotokos" {
				count++
			}
			if c.FeastName == "St Mary of Egypt" {
				count++
			}
			t.Log(c)
		}

		if count != 2 {
			t.Errorf("3/25/2018 should have The Annunciation and St. Mary of Egypt but doesn't.")
		}

		if day.FeastLevel != "Major feast Theotokos" {
			t.Errorf("3/25/2018 should have a feast level of 7 but doesn't.")
		}

		if strings.Contains(day.FastLevel, "Fish") {
			t.Errorf("3/25/2018 should be a fish day but isn't.")
		}
	})

	t.Run("Matins Gospel", func(t *testing.T) {
		// Veneration of the Cross - should include 7th Matins Gospel
		day := factory.NewDay(2018, 3, 11, nil)

		for _, r := range day.Readings {
			if r.Source == "7th Matins Gospel" {
				return
			}
		}

		t.Errorf("3/11/2018 should have the 7th Matins gospel but doesn't.")
	})

	t.Run("No memorial", func(t *testing.T) {
		// Memorial Saturday with no memorial readings
		// Should not have John 5.24-30
		day := factory.NewDay(2022, 3, 26, nil)

		for _, r := range day.Readings {
			if r.ShortDisplay == "John 5.24-30" {
				t.Errorf("3/26/2022 should not have John 5.24-30 but does.")
			}
		}
	})

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
}
