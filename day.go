package orthocal

import (
	"database/sql"
	"log"
)

type Day struct {
	Year    int
	Month   int
	Day     int
	PYear   *Year
	PDist   int
	Weekday int

	db *sql.DB
}

func NewDay(year, month, day int, useJulian bool, db *sql.DB) *Day {
	var self Day

	self.db = db
	self.Year, self.Month, self.Day = year, month, day
	pdist, pyear := ComputePaschaDistance(year, month, day)
	self.PDist = pdist
	self.PYear = NewYear(pyear, useJulian)
	self.Weekday = WeekDayFromPDist(self.PDist)

	return &self
}

func (self *Day) GetRecords() {
	var rows *sql.Rows
	var e error

	floatIndex := self.PYear.LookupFloatIndex(self.PDist)

	if floatIndex > 0 {
		rows, e = self.db.Query(
			`select title, feast_name, saint
			from days
			where pdist = $1 or pdist = $2
			or (month = $3 and day = $4)`, self.PDist, floatIndex, self.Month, self.Day)
	} else {
		rows, e = self.db.Query(
			`select title, feast_name, saint
			from days
			where pdist = $1
			or (month = $3 and day = $4)`, self.PDist, self.Month, self.Day)
	}

	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
	}

	for rows.Next() {
		var title, feastName, saint string
		rows.Scan(&title, &feastName, &saint)
		log.Printf("title = \"%s\", feast = \"%s\", saint = \"%s\"\n", title, feastName, saint)
	}
}
