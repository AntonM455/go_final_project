package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// afterNow compares only date (no time)
func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()

	if dy != ny {
		return dy > ny
	}
	if dm != nm {
		return dm > nm
	}
	return dd > nd
}

// NextDate calculates the next date according to the repeat rule
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// if the rule is empty, we delete
	if repeat == "" {
		return "", errors.New("no repeat rule: task will be deleted")
	}

	// parse the date dstart
	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid start date")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid number of days")
		}
		// move the dates until there are more now
		for {
			startDate = startDate.AddDate(0, 0, days)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	case "y":
		// move by years until no more now
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	default:
		return "", errors.New("unsupported repeat rule")
	}
}
