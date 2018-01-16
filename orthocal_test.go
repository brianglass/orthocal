package orthocal_test

import (
	"orthocal"
	"testing"
	"time"
)

var fixture_pascha = []time.Time{
	time.Date(2008, 4, 27, 0, 0, 0, 0, time.Local),
	time.Date(2009, 4, 19, 0, 0, 0, 0, time.Local),
	time.Date(2010, 4, 4, 0, 0, 0, 0, time.Local),
	time.Date(2011, 4, 24, 0, 0, 0, 0, time.Local),
}

func TestGregorianDateToJulianDay(t *testing.T) {
	julianDate := time.Date(2018, time.Month(1), 15, 0, 0, 0, 0, time.Local)
	actual := orthocal.GregorianDateToJulianDay(julianDate)
	expected := 2458134
	if actual != expected {
		t.Fatalf("GregorianDateToJulianDay should have returned %d but returned %d", expected, actual)
	}

	julianDate = time.Date(2000, time.Month(2), 29, 0, 0, 0, 0, time.Local)
	actual = orthocal.GregorianDateToJulianDay(julianDate)
	expected = 2451604
	if actual != expected {
		t.Fatalf("GregorianDateToJulianDay should have returned %d but returned %d", expected, actual)
	}
}

func TestCalculatePascha(t *testing.T) {
	for _, expectedTime := range fixture_pascha {
		pascha, e := orthocal.ComputeGregorianPascha(expectedTime.Year())
		if e != nil {
			t.Fatalf("CalculatePascha had an error: %#n", e)
		}
		if pascha != expectedTime {
			t.Fatalf("CalculatePascha should have returned %s but returned %s", expectedTime, pascha)
		}
	}
}
