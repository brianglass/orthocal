package orthocal_test

import (
	"orthocal"
	"testing"
	"time"
)

var fixture_gregorian_pascha = []time.Time{
	time.Date(2008, 4, 27, 0, 0, 0, 0, time.Local),
	time.Date(2009, 4, 19, 0, 0, 0, 0, time.Local),
	time.Date(2010, 4, 4, 0, 0, 0, 0, time.Local),
	time.Date(2011, 4, 24, 0, 0, 0, 0, time.Local),
}

var fixture_julian_pascha = []time.Time{
	time.Date(2008, 4, 14, 0, 0, 0, 0, time.Local),
	time.Date(2009, 4, 6, 0, 0, 0, 0, time.Local),
	time.Date(2010, 3, 22, 0, 0, 0, 0, time.Local),
	time.Date(2011, 4, 11, 0, 0, 0, 0, time.Local),
}

func TestGregorianDateToJulianDay(t *testing.T) {
	julianDate := time.Date(2018, time.Month(1), 15, 0, 0, 0, 0, time.Local)
	actual := orthocal.GregorianDateToJulianDay(julianDate)
	expected := 2458134
	if actual != expected {
		t.Fatalf("GregorianDateToJulianDay should have returned %d but returned %d", expected, actual)
	}

	julianDate = time.Date(2000, time.Month(5), 29, 0, 0, 0, 0, time.Local)
	actual = orthocal.GregorianDateToJulianDay(julianDate)
	expected = 2451694
	if actual != expected {
		t.Fatalf("GregorianDateToJulianDay should have returned %d but returned %d", expected, actual)
	}
}

func TestComputeGregorianPascha(t *testing.T) {
	for _, expectedTime := range fixture_gregorian_pascha {
		pascha, e := orthocal.ComputeGregorianPascha(expectedTime.Year())
		if e != nil {
			t.Errorf("CalculateGregorianPascha had an error: %#n", e)
		}
		if pascha != expectedTime {
			t.Errorf("CalculateGregorianPascha should have returned %s but returned %s", expectedTime, pascha)
		}
	}
}

func TestComputeGregorianPaschaInvalid(t *testing.T) {
	expectedTime := time.Date(2100, 5, 2, 0, 0, 0, 0, time.Local)
	_, e := orthocal.ComputeGregorianPascha(expectedTime.Year())
	if e == nil {
		t.Errorf("CalculateGregorianPacha should return an error when dates are out of range")
	}
}

func TestComputeJulianPascha(t *testing.T) {
	for _, expectedTime := range fixture_julian_pascha {
		month, day := orthocal.ComputeJulianPascha(expectedTime.Year())
		if time.Month(month) != expectedTime.Month() {
			t.Errorf("CalculateJulianPascha should have returned %d for month but returned %d", month, expectedTime.Month())
		}
		if day != expectedTime.Day() {
			t.Errorf("CalculateJulianPascha should have returned %d for day but returned %d", day, expectedTime.Day())
		}
	}
}

func TestConvertJulianToGregorian(t *testing.T) {
	expected := time.Date(2008, 4, 27, 0, 0, 0, 0, time.Local)
	actual, e := orthocal.JulianToGregorian(2008, 4, 14)
	if e != nil {
		t.Fatalf("ConvertJulianToGregory return error: %#n", e)
	}

	if expected != actual {
		t.Errorf("ConvertJulianToGregory should have returned %s but returned %d", expected, actual)
	}

	expected = time.Date(2011, 4, 24, 0, 0, 0, 0, time.Local)
	actual, e = orthocal.JulianToGregorian(2011, 4, 11)
	if e != nil {
		t.Fatalf("ConvertJulianToGregory return error: %#n", e)
	}

	if expected != actual {
		t.Errorf("ConvertJulianToGregory should have returned %s but returned %d", expected, actual)
	}
}

func TestConvertJulianToGregorianInvalid(t *testing.T) {
	_, e := orthocal.JulianToGregorian(2100, 4, 14)
	if e == nil {
		t.Errorf("ConvertJulianToGregory should return an error when dates are out of range")
	}

	_, e = orthocal.JulianToGregorian(1900, 4, 14)
	if e == nil {
		t.Errorf("ConvertJulianToGregory should return an error when dates are out of range")
	}
}

func TestJulianDatetoJulianDay(t *testing.T) {
	expected := 2455676
	actual := orthocal.JulianDateToJulianDay(2011, 4, 11)
	if actual != expected {
		t.Errorf("GregorianDateToJulianDay returned %d but should have returned %d", actual, expected)
	}
}

func TestComputJulianDayPascha(t *testing.T) {
	expected := 2455676
	actual := orthocal.ComputeJulianDayPascha(2011)
	if actual != expected {
		t.Errorf("ComputJulianDayPascha returned %d but should have returned %d", actual, expected)
	}
}

func TestComputePaschaDistance(t *testing.T) {
	date := time.Date(2018, 5, 9, 0, 0, 0, 0, time.Local)
	distance, year := orthocal.ComputePaschaDistance(date)
	expectedDistance, expectedYear := 31, 2018
	if distance != expectedDistance {
		t.Errorf("ComputePaschaDistance returned %d for the distance but should have returned %d", distance, expectedDistance)
	}
	if year != expectedYear {
		t.Errorf("ComputePaschaDistance returned %d for the year but should have returned %d", year, expectedYear)
	}

	date = time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local)
	distance, year = orthocal.ComputePaschaDistance(date)
	expectedDistance, expectedYear = 260, 2017
	if distance != expectedDistance {
		t.Errorf("ComputePaschaDistance returned %d for the distance but should have returned %d", distance, expectedDistance)
	}
	if year != expectedYear {
		t.Errorf("ComputePaschaDistance returned %d for the year but should have returned %d", year, expectedYear)
	}
}
