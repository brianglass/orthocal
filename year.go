package orthocal

type Year struct {
	Year int

	Pascha     int // Julian Day Number (JDN)
	NextPascha int

	// These measure the distance from Pascha (PDist)
	Finding              int
	Annunciation         int
	PeterAndPaul         int
	Beheading            int
	NativityTheotokos    int
	Elevation            int
	FathersSix           int
	FathersSeven         int
	DemetriusSaturday    int
	SynaxisUnmercenaries int
	Nativity             int
	Forefathers          int
	Theophany            int

	ExtraSundays int

	// This is the number of days after the Elevation?
	LucanJump int

	Reserves []int

	// unexported
	floats    []float
	noDaily   map[int]bool
	useJulian bool
}

type float struct {
	Index int
	PDist int
}

func NewYear(year int, useJulian bool) *Year {
	var self Year

	self.floats = make([]float, 0, 38)
	self.noDaily = make(map[int]bool)

	self.useJulian = useJulian
	self.Year = year
	self.Pascha = ComputePaschaJDN(year)
	self.NextPascha = ComputePaschaJDN(year + 1)
	self.computePDists()
	self.computeFloats()
	self.computeNoDailyReadings()
	self.computeReserves()

	return &self
}

func (self *Year) LookupFloatIndex(pdist int) int {
	// Since the stuff at the top is higher priority than the stuff at the
	// bottom, we do a linear search.
	for _, float := range self.floats {
		if float.PDist == pdist {
			return float.Index
		}
	}

	return 499
}

func (self *Year) HasNoDailyReadings(pdist int) bool {
	_, exists := self.noDaily[pdist]
	return exists
}

func (self *Year) dateToPDist(month, day, year int) int {
	if self.useJulian {
		// TODO: Need to test this and confirm it's valid
		return JulianDateToJDN(year, month, day) - self.Pascha
	} else {
		return GregorianDateToJDN(year, month, day) - self.Pascha
	}
}

// Compute the distance from Pascha for important feast days.
func (self *Year) computePDists() {
	var pdist, weekday int // for intermediate results

	self.Theophany = self.dateToPDist(1, 6, self.Year+1)
	self.Finding = self.dateToPDist(2, 24, self.Year)
	self.Annunciation = self.dateToPDist(3, 25, self.Year)
	self.PeterAndPaul = self.dateToPDist(6, 29, self.Year)

	// The Fathers of the Sixth Ecumenical Council falls on the Sunday nearest 7/16
	pdist = self.dateToPDist(7, 16, self.Year)
	weekday = WeekDayFromPDist(pdist)
	if weekday < Thursday {
		self.FathersSix = pdist - weekday
	} else {
		self.FathersSix = pdist + 7 - weekday
	}

	self.Beheading = self.dateToPDist(8, 29, self.Year)
	self.NativityTheotokos = self.dateToPDist(9, 8, self.Year)
	self.Elevation = self.dateToPDist(9, 14, self.Year)

	// The Fathers of the Seventh Ecumenical Council falls on the Sunday
	// following 10/11 or 10/11 itself if it is a Sunday.
	pdist = self.dateToPDist(10, 11, self.Year)
	weekday = WeekDayFromPDist(pdist)
	if weekday > Sunday {
		pdist += 7 - weekday
	}
	self.FathersSeven = pdist

	// Demetrius Saturday is the Saturday before 10/26
	pdist = self.dateToPDist(10, 26, self.Year)
	self.DemetriusSaturday = pdist - WeekDayFromPDist(pdist) - 1

	// The Synaxis of the Unmercenaries is the Sunday following 11/1
	pdist = self.dateToPDist(11, 1, self.Year)
	self.SynaxisUnmercenaries = pdist + 7 - WeekDayFromPDist(pdist)

	self.Nativity = self.dateToPDist(12, 25, self.Year)

	// Forefathers Sunday is the week before the week of Nativity
	weekday = WeekDayFromPDist(self.Nativity)
	self.Forefathers = self.Nativity - 14 + ((7 - weekday) % 7)

	// 168 - (Sunday after Elevation)
	self.LucanJump = 168 - (self.Elevation + 7 - WeekDayFromPDist(self.Elevation))
}

func (self *Year) addFloat(index, pdist int) {
	self.floats = append(self.floats, float{index, pdist})
}

func (self *Year) computeFloats() {
	// Order matters since we do a sequential search for the pdist values. The
	// stuff at the top has higher priority than the stuff at the bottom.

	self.addFloat(1001, self.FathersSix)
	self.addFloat(1002, self.FathersSeven)
	self.addFloat(1003, self.DemetriusSaturday)
	self.addFloat(1004, self.SynaxisUnmercenaries)

	// Floats around the Elevation of the Cross
	satBefore, sunBefore, satAfter, sunAfter := SurroundingWeekends(self.Elevation)
	if satBefore == self.NativityTheotokos {
		self.addFloat(1005, self.Elevation-1)
	} else {
		self.addFloat(1006, satBefore)
	}
	self.addFloat(1007, sunBefore)
	self.addFloat(1008, satAfter)
	self.addFloat(1009, sunAfter)
	self.addFloat(1010, self.Forefathers)

	// Floats around Nativity
	satBefore, sunBefore, satAfter, sunAfter = SurroundingWeekends(self.Nativity)
	switch self.Nativity - 1 {
	case satBefore:
		self.addFloat(1013, self.Nativity-2)
		self.addFloat(1012, sunBefore)
		self.addFloat(1015, self.Nativity-1)
	case sunBefore:
		self.addFloat(1013, self.Nativity-3)
		self.addFloat(1011, sunBefore)
		self.addFloat(1016, self.Nativity-1)
	default:
		self.addFloat(1014, self.Nativity-1)
		self.addFloat(1011, satBefore)
		self.addFloat(1012, sunBefore)
	}

	satBeforeTheophany, sunBeforeTheophany, satAfterTheophany, sunAfterTheophany := SurroundingWeekends(self.Theophany)
	switch WeekDayFromPDist(self.Nativity) {
	case Sunday:
		self.addFloat(1017, satAfter)
		self.addFloat(1020, self.Nativity+1)
		self.addFloat(1024, sunBeforeTheophany)
		self.addFloat(1026, self.Theophany-1)
	case Monday:
		self.addFloat(1017, satAfter)
		self.addFloat(1021, sunAfter)
		self.addFloat(1023, self.Theophany-5)
		self.addFloat(1026, self.Theophany-1)
	case Tuesday:
		self.addFloat(1019, satAfter)
		self.addFloat(1021, sunAfter)
		self.addFloat(1027, satBeforeTheophany)
		self.addFloat(1023, self.Theophany-5)
		self.addFloat(1025, self.Theophany-2)
	case Wednesday:
		self.addFloat(1019, satAfter)
		self.addFloat(1021, sunAfter)
		self.addFloat(1022, satBeforeTheophany)
		self.addFloat(1028, sunBeforeTheophany)
		self.addFloat(1025, self.Theophany-3)
	case Thursday, Friday:
		self.addFloat(1019, satAfter)
		self.addFloat(1021, sunAfter)
		self.addFloat(1022, satBeforeTheophany)
		self.addFloat(1024, sunBeforeTheophany)
		self.addFloat(1026, self.Theophany-1)
	case Saturday:
		self.addFloat(1018, self.Nativity+6)
		self.addFloat(1021, sunAfter)
		self.addFloat(1022, satBeforeTheophany)
		self.addFloat(1024, sunBeforeTheophany)
		self.addFloat(1026, self.Theophany-1)
	}
	self.addFloat(1029, satAfterTheophany)
	self.addFloat(1030, sunAfterTheophany)

	// New Martyrs of Russia (OCA) is the Sunday on or before 1/31
	martyrs := self.dateToPDist(1, 31, self.Year)
	weekday := WeekDayFromPDist(martyrs)
	if weekday != Sunday {
		// The Sunday before 1/31
		martyrs = martyrs - 7 + ((7 - weekday) % 7)
	}
	self.addFloat(1031, martyrs)

	// Floats around Annunciation
	switch WeekDayFromPDist(self.Annunciation) {
	case Saturday:
		self.addFloat(1032, self.Annunciation-1)
		self.addFloat(1033, self.Annunciation)
	case Sunday:
		self.addFloat(1034, self.Annunciation)
	case Monday:
		self.addFloat(1035, self.Annunciation)
	default:
		self.addFloat(1036, self.Annunciation-1)
		self.addFloat(1037, self.Annunciation)
	}
}

// assemble list of days on which daily readings are suppressed
func (self *Year) computeNoDailyReadings() {
	_, sunBefore, satAfter, sunAfter := SurroundingWeekends(self.Theophany)
	self.noDaily[sunBefore] = true
	self.noDaily[sunAfter] = true

	self.noDaily[self.Theophany-5] = true
	self.noDaily[self.Theophany-1] = true
	self.noDaily[self.Theophany] = true

	if satAfter == self.Theophany+1 {
		self.noDaily[self.Theophany+1] = true
	}

	self.noDaily[self.Forefathers] = true

	_, sunBefore, _, sunAfter = SurroundingWeekends(self.Nativity)
	self.noDaily[sunBefore] = true
	self.noDaily[self.Nativity-1] = true
	self.noDaily[self.Nativity] = true
	self.noDaily[self.Nativity+1] = true
	self.noDaily[sunAfter] = true

	if WeekDayFromPDist(self.Annunciation) == Saturday {
		self.noDaily[self.Annunciation] = true
	}
}

func (self *Year) computeReserves() {
	// TODO: store surrounding weekends in the struct
	_, _, _, sunAfter := SurroundingWeekends(self.Theophany)
	self.ExtraSundays = (self.NextPascha - self.Pascha - 84 - sunAfter) / 7

	if self.ExtraSundays > 0 {
		for i := self.Forefathers + self.LucanJump + 7; i <= 266; i += 7 {
			self.Reserves = append(self.Reserves, i)
		}
		remainder := self.ExtraSundays - len(self.Reserves)
		if remainder > 0 {
			for i := 175 - remainder*7; i < 169; i += 7 {
				self.Reserves = append(self.Reserves, i)
			}
		}
	}
}
