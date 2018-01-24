package orthocal

type Year struct {
	Year int

	Pascha int // Julian Day Number (JDN)

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

	Floats map[int]int

	// This is the number of days after the Elevation?
	LucanJump int

	// unexported
	useJulian bool
}

func NewYear(year int, useJulian bool) *Year {
	var self Year

	self.Floats = make(map[int]int)

	self.useJulian = useJulian
	self.Year = year
	self.Pascha = ComputePaschaJDN(year)
	self.computePDists()
	self.computeFloats()

	return &self
}

func (self *Year) dateToPDist(month, day int) int {
	if self.useJulian {
		// TODO: Need to test this and confirm it's valid
		return JulianDateToJDN(self.Year, month, day) - self.Pascha
	} else {
		return GregorianDateToJDN(self.Year, month, day) - self.Pascha
	}
}

// Compute the distance from Pascha for important feast days.
func (self *Year) computePDists() {
	var pdist, weekday int // for intermediate results

	self.Theophany = self.dateToPDist(1, 6)
	self.Finding = self.dateToPDist(2, 24)
	self.Annunciation = self.dateToPDist(3, 25)
	self.PeterAndPaul = self.dateToPDist(6, 29)

	// The Fathers of the Sixth Ecumenical Council falls on the Sunday nearest 7/16
	pdist = self.dateToPDist(7, 16)
	weekday = WeekDayFromPDist(pdist)
	if weekday < Thursday {
		self.FathersSix = pdist - weekday
	} else {
		self.FathersSix = pdist + 7 - weekday
	}

	self.Beheading = self.dateToPDist(8, 29)
	self.NativityTheotokos = self.dateToPDist(9, 8)
	self.Elevation = self.dateToPDist(9, 14)

	// The Fathers of the Seventh Ecumenical Council falls on the Sunday
	// following 10/11 or 10/11 itself if it is a Sunday.
	pdist = self.dateToPDist(10, 11)
	weekday = WeekDayFromPDist(pdist)
	if weekday > Sunday {
		pdist += 7 - weekday
	}
	self.FathersSeven = pdist

	// Demetrius Saturday is the Saturday before 10/26
	pdist = self.dateToPDist(10, 26)
	self.DemetriusSaturday = pdist - WeekDayFromPDist(pdist) - 1

	// The Synaxis of the Unmercenaries is the Sunday following 11/1
	pdist = self.dateToPDist(11, 1)
	self.SynaxisUnmercenaries = pdist + 7 - WeekDayFromPDist(pdist)

	self.Nativity = self.dateToPDist(12, 25)

	// Forefathers Sunday is the week before the week of Nativity
	weekday = WeekDayFromPDist(self.Nativity)
	self.Forefathers = self.Nativity - 14 + ((7 - weekday) % 7)

	// 168 - (Sunday after Elevation)
	self.LucanJump = 168 - (self.Elevation + 7 - WeekDayFromPDist(self.Elevation))
}

func (self *Year) computeFloats() {
	for i := 1001; i < 1038; i++ {
		self.Floats[i] = 499
	}

	self.Floats[1001] = self.FathersSix
	self.Floats[1002] = self.FathersSeven
	self.Floats[1003] = self.DemetriusSaturday
	self.Floats[1004] = self.SynaxisUnmercenaries

	// Floats around the Elevation of the Cross
	satBefore, sunBefore, satAfter, sunAfter := SurroundingWeekends(self.Elevation)
	if satBefore == self.NativityTheotokos {
		self.Floats[1005] = self.Elevation - 1
	} else {
		self.Floats[1006] = satBefore
	}
	self.Floats[1007] = sunBefore
	self.Floats[1008] = sunAfter
	self.Floats[1009] = satAfter
	self.Floats[1010] = sunAfter

	// Floats around Nativity
	satBefore, sunBefore, satAfter, sunAfter = SurroundingWeekends(self.Nativity)
	switch self.Nativity - 1 {
	case satBefore:
		self.Floats[1012] = sunBefore
		self.Floats[1013] = self.Nativity - 2
		self.Floats[1015] = self.Nativity - 1
	case sunBefore:
		self.Floats[1011] = satBefore
		self.Floats[1013] = self.Nativity - 3
		self.Floats[1016] = self.Nativity - 1
	default:
		self.Floats[1011] = satBefore
		self.Floats[1012] = sunBefore
		self.Floats[1014] = self.Nativity - 1
	}

	satBeforeTheophany, sunBeforeTheophany, satAfterTheophany, sunAfterTheophany := SurroundingWeekends(self.Theophany)
	switch WeekDayFromPDist(self.Nativity) {
	case Sunday:
		self.Floats[1017] = satAfter
		self.Floats[1020] = self.Nativity + 1
		self.Floats[1024] = sunBeforeTheophany
		self.Floats[1026] = self.Theophany - 1
	case Monday:
		self.Floats[1017] = satAfter
		self.Floats[1021] = sunAfter
		self.Floats[1023] = self.Theophany - 5
		self.Floats[1026] = self.Theophany - 1
	case Tuesday:
		self.Floats[1019] = satAfter
		self.Floats[1021] = sunAfter
		self.Floats[1027] = satBeforeTheophany
		self.Floats[1023] = self.Theophany - 5
		self.Floats[1025] = self.Theophany - 2
	case Wednesday:
		self.Floats[1019] = satAfter
		self.Floats[1021] = sunAfter
		self.Floats[1022] = satBeforeTheophany
		self.Floats[1028] = sunBeforeTheophany
		self.Floats[1025] = self.Theophany - 3
	case Thursday, Friday:
		self.Floats[1019] = satAfter
		self.Floats[1021] = sunAfter
		self.Floats[1022] = satBeforeTheophany
		self.Floats[1024] = sunBeforeTheophany
		self.Floats[1026] = self.Theophany - 1
	case Saturday:
		self.Floats[1018] = self.Nativity + 6
		self.Floats[1021] = sunAfter
		self.Floats[1022] = satBeforeTheophany
		self.Floats[1024] = sunBeforeTheophany
		self.Floats[1026] = self.Theophany - 1
	}
	self.Floats[1029] = satAfterTheophany
	self.Floats[1030] = sunAfterTheophany

	// Floats around Annunciation
}
