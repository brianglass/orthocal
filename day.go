package orthocal

import (
	"database/sql"
	"log"
)

type Day struct {
	PDist          int
	Year           int
	Month          int
	Day            int
	Weekday        int
	FeastLevel     string
	FastLevel      string
	FastException  string
	Commemorations []Commemoration

	db    *sql.DB
	pyear *Year
}

type Commemoration struct {
	Title     string
	Subtitle  string
	FeastName string
	SaintNote string
	Saint     string
}

func NewDay(year, month, day int, useJulian bool, db *sql.DB) *Day {
	var self Day

	self.db = db
	self.Year, self.Month, self.Day = year, month, day
	pdist, pyear := ComputePaschaDistance(year, month, day)
	self.PDist = pdist
	self.pyear = NewYear(pyear, useJulian)
	self.Weekday = WeekDayFromPDist(self.PDist)

	return &self
}

func (self *Day) GetRecords() {
	var rows *sql.Rows
	var e error

	floatIndex := self.pyear.LookupFloatIndex(self.PDist)

	if floatIndex > 0 {
		rows, e = self.db.Query(
			`select title, subtitle, feast_name, feast_level, saint_note, saint, fast, fast_exception
			from days
			where pdist = $1 or pdist = $2
			or (month = $3 and day = $4)`, self.PDist, floatIndex, self.Month, self.Day)
	} else {
		rows, e = self.db.Query(
			`select title, subtitle, feast_name, feast_level, saint_note, saint, fast, fast_exception
			from days
			from days
			where pdist = $1
			or (month = $3 and day = $4)`, self.PDist, self.Month, self.Day)
	}

	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
	}

	var overallFastLevel, overallFastException, overallFeastLevel int
	for rows.Next() {
		var title, subtitle, feastName, saintNote, saint string
		var feastLevel, fast, fastException int

		rows.Scan(&title, &subtitle, &feastName, &feastLevel, &saintNote, &saint, &fast, &fastException)
		c := Commemoration{title, subtitle, feastName, saintNote, saint}
		self.Commemorations = append(self.Commemorations, c)

		if feastLevel > overallFeastLevel {
			overallFeastLevel = feastLevel
		}
		if fast > overallFastLevel {
			overallFastLevel = fast
		}
		if fastException > overallFastException {
			overallFastException = fastException
		}

		self.FastLevel = FastLevels[overallFastLevel]
		self.FastException = FastExceptions[overallFastException]
		self.FeastLevel = FeastLevels[overallFeastLevel]
	}
}

/*
func (self *Day) GetReadings() []string {
	sql := `
		select readings.*, zachalos.zaDisplay as display, zachalos.zaSdisplay as sdisplay
		from readings left join zachalos
		on (zachalos.zaBook=readings.reBook and zachalos.zaNum=readings.reNum)
		where $conds
		order by reIndex`
}
*/
