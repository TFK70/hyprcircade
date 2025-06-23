package time

import (
	"testing"
)

func assertTod(t *testing.T, darkAt int, lightAt int, time int, expectedTod string) {
	tod, err := GetTimeOfTheDay(darkAt, lightAt, time)
	if err != nil {
		t.Fatalf("Unexpected error was thrown: %s", err.Error())
	}
	if tod != expectedTod {
		t.Fatalf(
			"Wrong time of day calculation:\nTime: %d\nDark at: %d\nLight at: %d\nExpected tod: %s\nReceived tod: %s",
			time,
			darkAt,
			lightAt,
			expectedTod,
			tod,
		)
	}
}

func TestGetTimeOfDay(t *testing.T) {
	darkAt := 20
	lightAt := 8

	time1 := 6
	time2 := 8
	time3 := 10
	time4 := 15
	time5 := 20
	time6 := 21
	time7 := 23
	time8 := 0

	assertTod(t, darkAt, lightAt, time1, DARK)
	assertTod(t, darkAt, lightAt, time2, LIGHT)
	assertTod(t, darkAt, lightAt, time3, LIGHT)
	assertTod(t, darkAt, lightAt, time4, LIGHT)
	assertTod(t, darkAt, lightAt, time5, DARK)
	assertTod(t, darkAt, lightAt, time6, DARK)
	assertTod(t, darkAt, lightAt, time7, DARK)
	assertTod(t, darkAt, lightAt, time8, DARK)
}
