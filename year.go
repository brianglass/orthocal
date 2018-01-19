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

	// This is the number of days after the Elevation
	LucanJump int

	// unexported
	useJulian bool
}

func NewYear(year int, useJulian bool) *Year {
	var self Year

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
}
