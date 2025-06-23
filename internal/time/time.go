package time

import (
	t "time"
)

const (
	DARK  = "dark"
	LIGHT = "light"
)

func GetCurrentTimeOfTheDay(darkAt int, lightAt int) (string, error) {
	now := t.Now()
	hour := now.Hour()

	tod, err := GetTimeOfTheDay(darkAt, lightAt, hour)
	if err != nil {
		return "", err
	}

	return tod, nil
}

func GetTimeOfTheDay(darkAt int, lightAt int, hour int) (string, error) {
	isDark := hour < lightAt || hour >= darkAt

	if isDark {
		return DARK, nil
	}

	return LIGHT, nil
}
