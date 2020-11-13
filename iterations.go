package tapfn

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bmheenan/taps"
)

// If you need a thread done by `neededBy`, being done by someone who plans by `cadence`, this returns the iteration
// they must plan it for
func iterRequired(neededBy string, cadence taps.Cadence) (string, error) {
	// TODO This arguably should be a different implementation than iterResulting. Current implementation is not strict:
	// it's possibly to land a thread later than the requested implementation, if at the end of the iteration. For now,
	// good enough; it's usually right
	return iterResulting(neededBy, cadence)
}

// If you have a thread that will be done by `doneBy`, this returns which iteration you can expect it in, with a cadence
// matching `cadence`
func iterResulting(doneBy string, cadence taps.Cadence) (string, error) {
	endDt, errED := getEndDate(doneBy)
	if errED != nil {
		return "", fmt.Errorf("Could not get end date of given iteration: %v", errED)
	}
	iter, errIt := iterContaining(endDt, cadence)
	if errIt != nil {
		return "", fmt.Errorf("Could not get iteration from end date %v: %v", endDt, errIt)
	}
	return iter, nil
}

func getEndDate(iter string) (time.Time, error) {
	yrly, errY := regexp.MatchString("^[0-9]{4}$", iter)
	if errY != nil {
		return time.Time{}, fmt.Errorf("Could not check if iteration was yearly: %v", errY)
	}
	if yrly {
		y, _ := strconv.Atoi(iter)
		return time.Date(y, time.Month(12), 31, 0, 0, 0, 0, time.UTC), nil
	}
	qrly, errQ := regexp.MatchString("^[0-9]{4} [Qq][1-4]$", iter)
	if errQ != nil {
		return time.Time{}, fmt.Errorf("Could not check if iteration was quarterly: %v", errQ)
	}
	if qrly {
		y, _ := strconv.Atoi(iter[0:4])
		q, _ := strconv.Atoi(iter[6:])
		var m, d int
		if q == 1 {
			m = 3
			d = 31
		} else if q == 2 {
			m = 6
			d = 30
		} else if q == 3 {
			m = 9
			d = 30
		} else {
			m = 12
			d = 31
		}
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
	}
	mnly, errM := regexp.MatchString("^[0-9]{4} [A-Za-z]{3}$", iter)
	if errM != nil {
		return time.Time{}, fmt.Errorf("Could not check if iteration was monthly: %v", errM)
	}
	if mnly {
		y, _ := strconv.Atoi(iter[0:4])
		mo := strings.ToLower(iter[5:])
		var d, m int
		if mo == "jan" {
			m = 1
			d = 31
		} else if mo == "feb" {
			m = 2
			if y%4 == 0 {
				d = 29
			} else {
				d = 28
			}
		} else if mo == "mar" {
			m = 3
			d = 31
		} else if mo == "apr" {
			m = 4
			d = 30
		} else if mo == "may" {
			m = 5
			d = 31
		} else if mo == "jun" {
			m = 6
			d = 30
		} else if mo == "jul" {
			m = 7
			d = 31
		} else if mo == "aug" {
			m = 8
			d = 31
		} else if mo == "sep" {
			m = 9
			d = 30
		} else if mo == "oct" {
			m = 10
			d = 31
		} else if mo == "nov" {
			m = 11
			d = 30
		} else if mo == "dec" {
			m = 12
			d = 31
		} else {
			return time.Time{}, fmt.Errorf("Could not parse to a month: %v", mo)
		}
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
	}
	bwly, errB := regexp.MatchString("^[0-9]{4} [Ww][0-9]{1,2}$", iter)
	if errB != nil {
		return time.Time{}, fmt.Errorf("Could not check if iteration was biweekly: %v", errB)
	}
	if bwly {
		// TODO implement biweekly iterations
		return time.Time{}, errors.New("Biweekly iterations not implemented")
	}
	return time.Time{}, fmt.Errorf("Could not understand the IterTiming of iteration: %v", iter)
}

func iterContaining(date time.Time, cadence taps.Cadence) (string, error) {
	y := date.Year()
	if cadence == taps.Yearly {
		return fmt.Sprintf("%v", y), nil
	} else if cadence == taps.Quarterly {
		if date.Month() <= 3 {
			return fmt.Sprintf("%v Q1", y), nil
		} else if date.Month() <= 6 {
			return fmt.Sprintf("%v Q2", y), nil
		} else if date.Month() <= 9 {
			return fmt.Sprintf("%v Q3", y), nil
		} else {
			return fmt.Sprintf("%v Q4", y), nil
		}
	} else if cadence == taps.Monthly {
		mos := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		return fmt.Sprintf("%v %v", y, mos[date.Month()-1]), nil
	}
	return "", errors.New("Biweekly iterations not implemented")
}
