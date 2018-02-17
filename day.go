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
	DoJump         bool

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

func NewDay(year, month, day int, useJulian bool, doJump bool, db *sql.DB) *Day {
	var self Day

	self.db = db
	self.Year, self.Month, self.Day = year, month, day
	pdist, pyear := ComputePaschaDistance(year, month, day)
	self.PDist = pdist
	self.pyear = NewYear(pyear, useJulian)
	self.Weekday = WeekDayFromPDist(self.PDist)
	self.DoJump = doJump

	self.getRecords()

	return &self
}

func (self *Day) getRecords() {
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
	var jump int

	_, _, _, sunAfter := SurroundingWeekends(self.pyear.Elevation)
	if self.PDist > sunAfter && self.DoJump {
		jump = self.pyear.LucanJump
	} else {
		jump = 0
	}

			if ($day['no_memorial']) {$nomem=" and reDesc != 'Departed'";} else {$nomem="";}
			if ($day['gday'] != 499)
			  { $conditions[]="(rePday = {$day['gday']} and reType = 'Gospel' $nomem)"; }
			if ($day['eday'] != 499)
			  { $conditions[]="(rePday = {$day['eday']} and reType = 'Epistle' $nomem)"; }
			$conditions[]="(rePday = {$day['pday']} and reType != 'Epistle' and reType !='Gospel')";
			if ($day['fday'] && $day['fday'] != 499)
			  { $conditions[]="(rePday = {$day['fday']})"; }
			if ($day['matins_gospel'])
			  { $mg = $day['matins_gospel']+700; $conditions[]="(rePday = $mg)"; }
			if ($day['no_matins_gospel']) {$x="and reType != 'Matins Gospel'";} else {$x="";}
			if ($day['no_paremias']) {$y="and reType != 'Vespers'";} else {$y="";}
		// no readings for leavetaking annunciation on non-liturgy day
			if ($day['month']==3 && $day['day']==26 && ($day['dow']==1 || $day['dow']==2 || $day['dow']==4))
			{$z="and reDesc != 'Theotokos'";} else {$z="";}
			  $conditions[]="((reMonth = {$day['menaion_month']} and reDay = {$day['menaion_day']}) $y $x $z)";
			if ($day['get_paremias'])
			  { $pa=getdate(mktime(0, 0, 0, $day['month'], $day['day']+1, $day['year']));
			    $conditions[]="(reMonth = {$pa['mon']} and reDay = {$pa['mday']} and reType = 'Vespers')"; }
		// make sql
			$conds = implode(" or ", $conditions);
			$q = "select readings.*, zachalos.zaDisplay as display, zachalos.zaSdisplay as sdisplay from readings left join zachalos on (zachalos.zaBook=readings.reBook and zachalos.zaNum=readings.reNum) where $conds order by reIndex";

	sql := `
		select readings.*, zachalos.zaDisplay as display, zachalos.zaSdisplay as sdisplay
		from readings left join zachalos
		on (zachalos.zaBook=readings.reBook and zachalos.zaNum=readings.reNum)
		where $conds
		order by reIndex`

	return nil
}
*/
