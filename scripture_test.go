package orthocal_test

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orthocal"
	"testing"
)

func TestScriptureLookup(t *testing.T) {
	tests := []string{
		"Matt 1.1-25",
		"Matt 4.25-5.13",
		"Matt 10.32-36, 11.1",
		"Matt 6.31-34, 7.9-11",
		"Matt 10.1, 5-8",
	}
	db, e := sql.Open("sqlite3", "kjv.db")
	if e != nil {
		t.Errorf("Got error opening database: %#n.", e)
	}

	bible := orthocal.NewBible(db)

	for _, test := range tests {
		passage := bible.Lookup(test)
		for _, verse := range passage {
			t.Logf("%d:%d", verse.Chapter, verse.Verse)
		}
		t.Logf("")
	}

	t.Fail()
}
