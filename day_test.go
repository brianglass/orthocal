package orthocal_test

import (
	"database/sql"
	// "encoding/json"
	"github.com/brianglass/orthocal"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	// "time"
)

func TestDay(t *testing.T) {
	db, e := sql.Open("sqlite3", "oca_calendar.db")
	if e != nil {
		t.Errorf("Got error opening database: %#v.", e)
	}

	bibledb, e := sql.Open("sqlite3", "kjv.db")
	if e != nil {
		t.Errorf("Got error opening database: %#v.", e)
	}
	bible := orthocal.NewBible(bibledb)

	factory := orthocal.NewDayFactory(false, true, db)

	// Sunday of the Publican and Pharisee
	/*
		day := factory.NewDay(2018, 1, 28, bible)
		actual, _ := json.marshalindent(day, "", "\t")
		t.errorf("%s", actual)
	*/

	t.Run("Annunciation", func(t *testing.T) {
		day := factory.NewDay(2018, 3, 25, nil)

		count := 0
		for _, f := range day.Feasts {
			if f == "Annunciation Most Holy Theotokos" {
				count++
			}
			if f == "St Mary of Egypt" {
				count++
			}
		}

		if count != 2 {
			t.Errorf("3/25/2018 should have The Annunciation and St. Mary of Egypt but doesn't.")
		}

		if day.FeastLevel != 7 {
			t.Errorf("3/25/2018 should have a feast level of 7 but doesn't.")
		}

		if day.FastLevel == 4 {
			t.Errorf("3/25/2018 should be a fish day but isn't.")
		}

		if len(day.Readings) != 12 {
			t.Errorf("3/25/2018 should have 12 scripture readings but has %d.", len(day.Readings))
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

	t.Run("Scriptures", func(t *testing.T) {
		// Cheesefare Sunday
		day := factory.NewDay(2018, 2, 18, bible)

		if len(day.Readings) != 3 {
			t.Errorf("2/18/2018 should have 3 readings but has %d.", len(day.Readings))
		}

		if len(day.Readings[0].Passage) != 12 {
			t.Errorf("2/18/2018's first reading should be 12 verses long but is %d.", len(day.Readings[0].Passage))
		}

		if len(day.Readings[1].Passage) != 8 {
			t.Errorf("2/18/2018's second reading should be 8 verses long but is %d.", len(day.Readings[1].Passage))
		}

		if len(day.Readings[2].Passage) != 8 {
			t.Errorf("2/18/2018's third reading should be 8 verses long but is %d.", len(day.Readings[2].Passage))
		}
	})

	t.Run("Paremias", func(t *testing.T) {
		day := factory.NewDay(2018, 3, 8, bible)

		if len(day.Readings) != 6 {
			t.Errorf("3/8/2018 should have 6 readings but has %d.", len(day.Readings))
		}
	})

	t.Run("Sebaste", func(t *testing.T) {
		day := factory.NewDay(2018, 3, 9, bible)

		if len(day.Readings) != 6 {
			t.Errorf("3/9/2018 should have 6 readings but has %d.", len(day.Readings))
		}

		if day.Readings[0].Source != "Matins Gospel" {
			t.Errorf("3/9/2018 should have matins gospel first but doesn't.")
		}

		set := make(map[string]bool)
		for _, r := range day.Readings {
			set[r.ShortDisplay] = true
		}
		if len(set) != 6 {
			t.Errorf("3/9/2018 should not have duplicate readings.")
		}
	})

	t.Run("Tone", func(t *testing.T) {
		testCases := []struct {
			day  *orthocal.Day
			tone int
		}{
			{factory.NewDay(2018, 4, 12, nil), 0},
			{factory.NewDay(2018, 4, 17, nil), 1},
			{factory.NewDay(2018, 2, 6, nil), 2},
			{factory.NewDay(2019, 1, 23, nil), 1},
			{factory.NewDay(2019, 6, 21, nil), 7},
		}

		for _, tc := range testCases {
			t.Run("Tone", func(t *testing.T) {
				if tc.day.Tone != tc.tone {
					t.Errorf("%d/%d/%d should have tone %d but has tone %d.", tc.day.Month, tc.day.Day, tc.day.Year, tc.tone, tc.day.Tone)
				}
			})
		}
	})

	t.Run("Apostles Fast", func(t *testing.T) {
		testCases := []struct {
			day       *orthocal.Day
			fast      int
			exception int
		}{
			{factory.NewDay(2018, 6, 3, nil), 0, 0},
			{factory.NewDay(2018, 6, 4, nil), 3, 0},
			{factory.NewDay(2018, 6, 12, nil), 3, 1},
			{factory.NewDay(2018, 6, 14, nil), 3, 1},
			{factory.NewDay(2018, 6, 16, nil), 3, 2},
			{factory.NewDay(2018, 6, 17, nil), 3, 2},
			{factory.NewDay(2018, 6, 28, nil), 3, 1},
			{factory.NewDay(2018, 6, 29, nil), 1, 2},
			{factory.NewDay(2018, 6, 30, nil), 0, 0},
		}

		for _, tc := range testCases {
			t.Run("Tone", func(t *testing.T) {
				if tc.day.FastLevel != tc.fast {
					t.Errorf("%d/%d/%d should have fast level %d but has %d.", tc.day.Month, tc.day.Day, tc.day.Year, tc.fast, tc.day.FastLevel)
				}
				if tc.day.FastException != tc.exception {
					t.Errorf("%d/%d/%d should have fast exception %d but has %d.", tc.day.Month, tc.day.Day, tc.day.Year, tc.exception, tc.day.FastException)
				}
			})
		}
	})

	t.Run("Fast Free", func(t *testing.T) {
		testCases := []struct {
			day  *orthocal.Day
			fast int
			desc string
		}{
			{factory.NewDay(2018, 12, 26, nil), 0, "No Fast"},
			{factory.NewDay(2018, 12, 28, nil), 0, "No Fast"},
			{factory.NewDay(2019, 1, 2, nil), 0, "No Fast"},
			{factory.NewDay(2019, 1, 4, nil), 0, "No Fast"},
		}

		for _, tc := range testCases {
			t.Run("Tone", func(t *testing.T) {
				if tc.day.FastLevel != tc.fast {
					t.Errorf("%d/%d/%d should have fast level %d but has %d.", tc.day.Month, tc.day.Day, tc.day.Year, tc.fast, tc.day.FastLevel)
				}
				if tc.day.FastLevelDesc != tc.desc {
					t.Errorf("%d/%d/%d should have fast description \"%s\" but has \"%s\".", tc.day.Month, tc.day.Day, tc.day.Year, tc.desc, tc.day.FastLevelDesc)
				}
			})
		}
	})

	/*
		// today := time.Now()
		today := time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)
		for {
			today = today.AddDate(0, 0, 1)
			day := factory.NewDay(today.Year(), int(today.Month()), today.Day(), nil)
			if day.FastException == 11 && day.FastLevel != 0 {
				actual, _ := json.MarshalIndent(day, "", "\t")
				t.Errorf("%s", actual)
				break
			}
		}
	*/
}
