package orthocal

import (
	"database/sql"
	"log"
	"time"
)

type Day struct {
	PDist          int
	JDN            int
	Year           int
	Month          int
	Day            int
	Weekday        int
	FeastLevel     string
	FastLevel      string
	FastException  string
	Commemorations []Commemoration
	Readings       []Reading
	DoJump         bool

	db    *sql.DB
	pyear *Year
	bible *Bible
}

type Commemoration struct {
	Title     string
	Subtitle  string
	FeastName string
	SaintNote string
	Saint     string
}

type Reading struct {
	Source       string
	Book         string
	Description  string
	Display      string
	ShortDisplay string
	Passage      Passage
}

func NewDay(year, month, day int, useJulian bool, doJump bool, db *sql.DB, bible *Bible) *Day {
	var self Day

	self.db = db
	self.bible = bible

	// time.Date automatically wraps dates that are invalid to the next month.
	// e.g. April 31 -> May 1
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	self.Year, self.Month, self.Day = date.Year(), int(date.Month()), date.Day()

	pdist, pyear := ComputePaschaDistance(year, month, day)
	if useJulian {
		self.JDN = JulianDateToJDN(year, month, day)
	} else {
		self.JDN = GregorianDateToJDN(year, month, day)
	}
	self.PDist = pdist
	self.pyear = NewYear(pyear, useJulian)
	self.Weekday = WeekDayFromPDist(self.PDist)
	self.DoJump = doJump

	self.getCommemorations()
	self.getReadings()

	return &self
}

func (self *Day) getCommemorations() {
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
			where pdist = $1
			or (month = $3 and day = $4)`, self.PDist, self.Month, self.Day)
	}
	defer rows.Close()

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
func (self *Day) matinsGospel int {
	if self.Weekday != Sunday {
		return 0
	}

	if -8 < self.PDist < 50 {
		return 0
	}

	if self.FeastLevel < 7 {
	}

// calculate sunday matins gospel
	$no_matins_gospel=false;
	$matins_gospel=0;
	if ($dow==0) {
		if ($pday > -8 && $pday < 50) {
			$no_matins_gospel=true;
		} elseif ($feast_level < 7) {
			$no_matins_gospel=true;
			$mbase = $pbase - 49;
			$x = $mbase % 77;
			if ($x==0) {$x=77;}
			$matins_gospel = floor($x/7);
		}
	}
}
*/

func (self *Day) getReadings() {
	var gPDist, ePDist int
	var jump int

	// Compute the Lucan jump
	_, _, _, sunAfter := SurroundingWeekends(self.pyear.Elevation)
	if self.PDist > sunAfter && self.DoJump {
		jump = self.pyear.LucanJump
	} else {
		jump = 0
	}

	// Compute the adjusted pdists for epistle and gospel
	if self.pyear.HasNoDailyReadings(self.PDist) {
		gPDist, ePDist = 499, 499
	} else {
		limit := 272

		// Compute adjusted pdist for the epistle
		if self.PDist == 252 {
			ePDist = self.pyear.Forefathers
		} else if self.PDist > limit {
			ePDist = self.JDN - self.pyear.NextPascha
		} else {
			ePDist = self.PDist
		}

		if WeekDayFromPDist(self.pyear.Theophany) < Tuesday {
			limit = 279
		}

		// Compute adjusted pdist for the Gospel
		_, _, _, sunAfter := SurroundingWeekends(self.pyear.Theophany)
		if self.PDist == 245-self.pyear.LucanJump {
			gPDist = self.pyear.Forefathers + self.pyear.LucanJump
		} else if self.PDist > sunAfter && self.Weekday == Sunday && self.pyear.ExtraSundays > 1 {
			i := (self.PDist - sunAfter) / 7
			gPDist = self.pyear.Reserves[i-1]
		} else if self.PDist+jump > limit {
			// Theophany stepback
			gPDist = self.JDN - self.pyear.NextPascha
		} else {
			gPDist = self.PDist + jump
		}
	}

	// TODO: handle days with no memorials

	/*
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
	*/

	// Timings using Prepare instead of Query proved that the time saved on a
	// month of days was around a couple milliseconds and not worth the added
	// complexity.
	rows, e := self.db.Query(`
		select source, r.desc, r.book, display, sdisplay
		from readings r left join pericopes p
		on (r.book=p.book and r.pericope=p.pericope)
		where
			(pdist = $1 and source = 'Gospel')
			or (pdist = $2 and source = 'Epistle')
			or (pdist = $3 and source != 'Epistle' and source != 'Gospel')
			or (pdist = $4)
		order by ordering`, gPDist, ePDist, self.PDist, self.pyear.LookupFloatIndex(self.PDist))
	defer rows.Close()

	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
	}

	for rows.Next() {
		var reading Reading
		rows.Scan(&reading.Source, &reading.Description, &reading.Book, &reading.Display, &reading.ShortDisplay)
		if self.bible != nil {
			passage := self.bible.Lookup(reading.ShortDisplay)
			if passage != nil {
				reading.Passage = passage
			}
		}
		self.Readings = append(self.Readings, reading)
	}
}
