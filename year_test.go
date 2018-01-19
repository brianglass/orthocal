package orthocal_test

import (
	"orthocal"
	"testing"
)

func TestComputePDists(t *testing.T) {
	year := orthocal.NewYear(2018, false)
	pascha := orthocal.GregorianDateToJDN(2018, 4, 8)

	if pascha != year.Pascha {
		t.Errorf("Got incorrect date for Pascha: %s. Should be %s", year.Pascha, pascha)
	}

	theophany := orthocal.GregorianDateToJDN(2018, 1, 6) - pascha
	if theophany != year.Theophany {
		t.Errorf("Got incorrect date for Theophany: %d. Should be %d", year.Theophany, theophany)
	}
	finding := orthocal.GregorianDateToJDN(2018, 2, 24) - pascha
	if finding != year.Finding {
		t.Errorf("Got incorrect date for finding: %d. Should be %d", year.Finding, finding)
	}
	annunciation := orthocal.GregorianDateToJDN(2018, 3, 25) - pascha
	if annunciation != year.Annunciation {
		t.Errorf("Got incorrect date for annunciation: %d. Should be %d", year.Annunciation, annunciation)
	}
	peterandpaul := orthocal.GregorianDateToJDN(2018, 6, 29) - pascha
	if peterandpaul != year.PeterAndPaul {
		t.Errorf("Got incorrect date for Peter and Paul: %d. Should be %d", year.PeterAndPaul, peterandpaul)
	}
	fatherssix := orthocal.GregorianDateToJDN(2018, 7, 15) - pascha
	if fatherssix != year.FathersSix {
		t.Errorf("Got incorrect date for the Fathers of the first 6 councils: %d. Should be %d", year.FathersSix, fatherssix)
	}
	beheading := orthocal.GregorianDateToJDN(2018, 8, 29) - pascha
	if beheading != year.Beheading {
		t.Errorf("Got incorrect date for the Beheading: %d. Should be %d", year.Beheading, beheading)
	}
	nativitytheotokos := orthocal.GregorianDateToJDN(2018, 9, 8) - pascha
	if nativitytheotokos != year.NativityTheotokos {
		t.Errorf("Got incorrect date for the Nativity of the Theotokos: %d. Should be %d", year.NativityTheotokos, nativitytheotokos)
	}
	elevation := orthocal.GregorianDateToJDN(2018, 9, 14) - pascha
	if elevation != year.Elevation {
		t.Errorf("Got incorrect date for the Elevation of the Cross: %d. Should be %d", year.Elevation, elevation)
	}
	fathersseven := orthocal.GregorianDateToJDN(2018, 10, 14) - pascha
	if fathersseven != year.FathersSeven {
		t.Errorf("Got incorrect date for the Fathers of the seventh council: %d. Should be %d", year.FathersSeven, fathersseven)
	}
	demetrius := orthocal.GregorianDateToJDN(2018, 10, 20) - pascha
	if demetrius != year.DemetriusSaturday {
		t.Errorf("Got incorrect date for Demetrius Saturday: %d. Should be %d", year.DemetriusSaturday, demetrius)
	}
	forefathers := orthocal.GregorianDateToJDN(2018, 12, 16) - pascha
	if forefathers != year.Forefathers {
		t.Errorf("Got incorrect date for the Sunday of the Forefathers: %d. Should be %d", year.Forefathers, forefathers)
	}
	nativity := orthocal.GregorianDateToJDN(2018, 12, 25) - pascha
	if nativity != year.Nativity {
		t.Errorf("Got incorrect date for Nativity: %d. Should be %d", year.Nativity, nativity)
	}

	// TODO: Confirm this is actually working
	lucanjump := 7
	if lucanjump != year.LucanJump {
		t.Errorf("Got incorrect date for the Lucan jump: %d. Should be %d", year.LucanJump, lucanjump)
	}
}

func TestComputePDistsSixth(t *testing.T) {
	year := orthocal.NewYear(2016, false)
	pascha := orthocal.GregorianDateToJDN(2016, 5, 1)

	if pascha != year.Pascha {
		t.Errorf("Got incorrect date for Pascha: %s. Should be %s", year.Pascha, pascha)
	}

	fatherssix := orthocal.GregorianDateToJDN(2016, 7, 17) - pascha
	if fatherssix != year.FathersSix {
		t.Errorf("Got incorrect date for the Fathers of the first 6 councils: %d. Should be %d", year.FathersSix, fatherssix)
	}
}
