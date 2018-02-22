package orthocal_test

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"orthocal"
	"testing"
)

func TestScriptureLookup(t *testing.T) {
	db, e := sql.Open("sqlite3", "kjv.db")
	if e != nil {
		t.Errorf("Got error opening database: %#n.", e)
	}
	bible := orthocal.NewBible(db)

	testCases := []struct {
		reference string
		count     int
	}{
		{"Matt 1.1-25", 25},
		{"Matt 4.25-5.13", 14},
		{"Matt 10.32-36, 11.1", 6},
		{"Matt 6.31-34, 7.9-11", 7},
		{"Matt 10.1, 5-8", 5},
	}

	for _, tc := range testCases {
		message := fmt.Sprintf("Scripture Lookup")
		t.Run(message, func(t *testing.T) {
			passage := bible.Lookup(tc.reference)
			if len(passage) != tc.count {
				t.Logf("%s should return %d verses but returned %d verses.", tc.reference, tc.count, len(passage))
			}
		})
	}
}
