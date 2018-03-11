package orthocal

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Day struct {
	PDist          int       `json:"pascha_distance"`
	JDN            int       `json:"julian_day_number"`
	Year           int       `json:"year"`
	Month          int       `json:"month"`
	Day            int       `json:"day"`
	Weekday        int       `json:"weekday"`
	Titles         []string  `json:"titles"`
	FeastLevel     int       `json:"feast_level"`
	FeastLevelDesc string    `json:"feast_level_description"`
	Feasts         []string  `json:"feasts"`
	FastLevel      string    `json:"fast_level"`
	FastException  string    `json:"fast_exception"`
	Saints         []string  `json:"saints"`
	ServiceNotes   []string  `json:"service_notes"`
	Readings       []Reading `json:"readings"`

	pyear *Year
}

type Reading struct {
	Source       string  `json:"source"`
	Book         string  `json:"book"`
	Description  string  `json:"description"`
	Display      string  `json:"display"`
	ShortDisplay string  `json:"short_display"`
	Passage      Passage `json:"passage"`
}

func (self *Day) HasNoMemorial() bool {
	return ((self.PDist == -36 || self.PDist == -29 || self.PDist == -22) &&
		(self.Month == 3) &&
		(self.Day == 9 || self.Day == 24 || self.Day == 25 || self.Day == 26))
}

type DayFactory struct {
	db        *sql.DB
	useJulian bool
	doJump    bool
	years     sync.Map
}

func NewDayFactory(useJulian bool, doJump bool, db *sql.DB) *DayFactory {
	var self DayFactory
	self.db = db
	self.useJulian = useJulian
	self.doJump = doJump
	return &self
}

func (self *DayFactory) NewDay(year, month, day int, bible *Bible) *Day {
	var d Day

	// time.Date automatically wraps dates that are invalid to the next month.
	// e.g. April 31 -> May 1
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	d.Year, d.Month, d.Day = date.Year(), int(date.Month()), date.Day()

	pdist, pyear := ComputePaschaDistance(year, month, day)
	if self.useJulian {
		d.JDN = JulianDateToJDN(year, month, day)
	} else {
		d.JDN = GregorianDateToJDN(year, month, day)
	}
	d.PDist = pdist
	d.Weekday = WeekDayFromPDist(d.PDist)

	// Cache years in a thread-safe way
	if y, ok := self.years.Load(pyear); ok {
		d.pyear = y.(*Year)
	} else {
		d.pyear = NewYear(pyear, self.useJulian)
		self.years.Store(pyear, d.pyear)
	}

	self.addCommemorations(&d)
	self.addReadings(&d, bible)

	return &d
}

func (self *DayFactory) addCommemorations(day *Day) {
	var rows *sql.Rows
	var e error

	floatIndex := day.pyear.LookupFloatIndex(day.PDist)

	if floatIndex != 0 && floatIndex != 499 {
		rows, e = self.db.Query(
			`select title, subtitle, feast_name, feast_level, service_note, saint, fast, fast_exception
			from days
			where pdist = $1 or pdist = $2
			or (month = $3 and day = $4)`, day.PDist, floatIndex, day.Month, day.Day)
	} else {
		rows, e = self.db.Query(
			`select title, subtitle, feast_name, feast_level, service_note, saint, fast, fast_exception
			from days
			where pdist = $1
			or (month = $3 and day = $4)`, day.PDist, day.Month, day.Day)
	}

	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
		return
	}
	defer rows.Close()

	overallFastLevel, overallFastException, overallFeastLevel := 0, 0, -2
	for rows.Next() {
		var title, subtitle, feastName, serviceNote, saint string
		var feastLevel, fast, fastException int

		rows.Scan(&title, &subtitle, &feastName, &feastLevel, &serviceNote, &saint, &fast, &fastException)

		if len(subtitle) > 0 {
			title = fmt.Sprintf("%s: %s", title, subtitle)
		}

		if len(title) > 0 {
			day.Titles = append(day.Titles, title)
		}
		if len(saint) > 0 {
			day.Saints = append(day.Saints, saint)
		}
		if len(feastName) > 0 {
			day.Feasts = append(day.Feasts, feastName)
		}
		if len(serviceNote) > 0 {
			day.ServiceNotes = append(day.ServiceNotes, serviceNote)
		}

		// Composite values
		if feastLevel > overallFeastLevel {
			overallFeastLevel = feastLevel
		}
		if fast > overallFastLevel {
			overallFastLevel = fast
		}
		if fastException > overallFastException {
			overallFastException = fastException
		}

		day.FastLevel = FastLevels[overallFastLevel]
		day.FastException = FastExceptions[overallFastException]
		day.FeastLevel = overallFeastLevel
		day.FeastLevelDesc = FeastLevels[overallFeastLevel]
	}
}

func (self *DayFactory) addReadings(day *Day, bible *Bible) {
	var conditionals []string

	ePDist, gPDist := self.getAdjustedPDists(day)

	var departed string
	if day.HasNoMemorial() {
		departed = "and r.desc != 'Departed'"
	}

	// Conditional for floats
	floatIndex := day.pyear.LookupFloatIndex(day.PDist)
	if floatIndex != 499 {
		conditionals = append(conditionals, fmt.Sprintf("or (pdist = %d)", floatIndex))
	}

	// Conditional for Matins Gospel
	hasMatinsGospel, matinsGospel := self.matinsGospel(day)
	if matinsGospel != 0 {
		conditionals = append(conditionals, fmt.Sprintf("or (pdist = %d)", matinsGospel+700))
	}

	// Conditional for Paremias
	if day.pyear.HasParemias(day.PDist) {
		date := time.Date(day.Year, time.Month(day.Month), day.Day+1, 0, 0, 0, 0, time.Local)
		paremias := fmt.Sprintf("or (r.month = %d and r.day = %d and source = 'Vespers')", date.Month(), date.Day())
		conditionals = append(conditionals, paremias)
	}

	// Build Conditional for Month/Day (i.e. non-pdist)
	var m string
	if !hasMatinsGospel {
		m = "and r.source != 'Matins Gospel'"
	}
	var p string
	if day.pyear.HasNoParemias(day.PDist) {
		p = "and r.source != 'Vespers'"
	}
	var a string
	if day.Month == 3 && day.Day == 26 && (day.Weekday == Monday || day.Weekday == Tuesday || day.Weekday == Thursday) {
		// no readings for leavetaking annunciation on non-liturgy day
		a = "and r.desc != 'Theotokos'"
	}
	dates := fmt.Sprintf("or (r.month = %d and r.day = %d %s %s %s)", day.Month, day.Day, m, p, a)
	conditionals = append(conditionals, dates)

	// TODO: Handle arbitrary exceptions

	// Since no user provided strings are being used, it is safe to use
	// string interpolation to build the SQL.
	query := `
		select source, r.desc, p.book, display, sdisplay
		from readings r left join pericopes p
		on (r.book=p.book and r.pericope=p.pericope)
		where
			   (pdist = %d and source = 'Gospel' %s)
			or (pdist = %d and source = 'Epistle' %s)
			or (pdist = %d and source != 'Epistle' and source != 'Gospel')
			%s
		order by ordering`

	query = fmt.Sprintf(query, gPDist, departed, ePDist, departed, day.PDist, strings.Join(conditionals, " "))

	rows, e := self.db.Query(query)
	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
		return
	}
	defer rows.Close()

	// Fetch all the readings
	for rows.Next() {
		var reading Reading
		rows.Scan(&reading.Source, &reading.Description, &reading.Book, &reading.Display, &reading.ShortDisplay)
		if bible != nil {
			passage := bible.Lookup(reading.ShortDisplay)
			if passage != nil {
				reading.Passage = passage
			}
		}
		day.Readings = append(day.Readings, reading)
	}

	// Move Lenten Matins Gospel to the top
	if day.PDist > -42 && day.PDist < -7 && day.FeastLevel < 7 {
		for i, reading := range day.Readings {
			if reading.Source == "Matins Gospel" {
				// Remove the matins gospel from the slice
				x := append(day.Readings[:i], day.Readings[i+1:]...)
				// prepend the matins gospel to the slice
				day.Readings = append([]Reading{reading}, x...)
				break
			}
		}
	}
}

func (self *DayFactory) matinsGospel(day *Day) (bool, int) {
	if day.Weekday == Sunday {
		if day.PDist > -8 && day.PDist < 50 {
			return false, 0
		} else if day.FeastLevel < 7 {
			pbase := day.PDist
			if pbase < 0 {
				pbase = day.JDN - day.pyear.PreviousPascha
			}

			x := (pbase - 49) % 77
			if x == 0 {
				x = 77
			}

			return false, x / 7
		}
	}

	return true, 0
}

func (self *DayFactory) getAdjustedPDists(day *Day) (int, int) {
	var gPDist, ePDist int
	var jump int

	// Compute the Lucan jump
	_, _, _, sunAfter := SurroundingWeekends(day.pyear.Elevation)
	if day.PDist > sunAfter && self.doJump {
		jump = day.pyear.LucanJump
	}

	// Compute the adjusted pdists for epistle and gospel
	if day.pyear.HasNoDailyReadings(day.PDist) {
		gPDist, ePDist = 499, 499
	} else {
		limit := 272

		// Compute adjusted pdist for the epistle
		if day.PDist == 252 {
			ePDist = day.pyear.Forefathers
		} else if day.PDist > limit {
			ePDist = day.JDN - day.pyear.NextPascha
		} else {
			ePDist = day.PDist
		}

		if WeekDayFromPDist(day.pyear.Theophany) < Tuesday {
			limit = 279
		}

		// Compute adjusted pdist for the Gospel
		_, _, _, sunAfter := SurroundingWeekends(day.pyear.Theophany)
		if day.PDist == 245-day.pyear.LucanJump {
			gPDist = day.pyear.Forefathers + day.pyear.LucanJump
		} else if day.PDist > sunAfter && day.Weekday == Sunday && day.pyear.ExtraSundays > 1 {
			i := (day.PDist - sunAfter) / 7
			gPDist = day.pyear.Reserves[i-1]
		} else if day.PDist+jump > limit {
			// Theophany stepback
			gPDist = day.JDN - day.pyear.NextPascha
		} else {
			gPDist = day.PDist + jump
		}
	}

	return ePDist, gPDist
}
