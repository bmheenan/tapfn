package tapfn

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) iterCurrent(cadence taps.Cadence) string {
	var d time.Time
	if cn.timeOverride == (time.Time{}) {
		d = time.Now()
	} else {
		d = cn.timeOverride
	}
	return iterContaining(d, cadence)
}

func (cn *cnTapdb) itersBetweenCurrentAnd(end string) []string {
	cadence, valid := iterCadence(end)
	if !valid {
		panic(fmt.Sprintf("Could not get cadence of %v: invalid iteration", end))
	}
	now := cn.iterCurrent(cadence)
	var min, max string
	switch {
	case now == end:
		return []string{now}
	case now < end:
		max = end
		min = now
	case now > end:
		max = now
		min = end
	}
	iters := []string{min}
	for iters[len(iters)-1] != max {
		iters = append(iters, iterNext(iters[len(iters)-1]))
	}
	return iters
}

func (cn *cnTapdb) itersAddStd(in []string, cadence taps.Cadence) (out []string, err error) {
	current := cn.iterCurrent(cadence)
	if !strIn(current, in) {
		in = append(in, current)
	}
	next := iterNext(current)
	if !strIn(next, in) {
		in = append(in, next)
	}
	sort.Strings(in)
	if !strIn("Inbox", in) {
		in = append([]string{"Inbox"}, in...)
	}
	if !strIn("Backlog", in) {
		in = append(in, "Backlog")
	}
	out = in
	return
}

func strIn(s string, a []string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

// If you need a thread done by `neededBy`, being done by someone who plans by `cadence`, this returns the iteration
// they must plan it for
func iterRequired(neededBy string, cadence taps.Cadence) string {
	// TODO This arguably should be a different implementation than iterResulting. Current implementation is not strict:
	// it's possibly to land a thread later than the requested implementation, if at the end of the iteration. For now,
	// good enough; it's usually right
	return iterResulting(neededBy, cadence)
}

// If you have a thread that will be done by `doneBy`, this returns which iteration you can expect it in, with a cadence
// matching `cadence`
func iterResulting(doneBy string, cadence taps.Cadence) string {
	if doneBy == "Inbox" {
		return "Inbox"
	}
	if doneBy == "Backlog" {
		return "Backlog"
	}
	endDt := iterEndDate(doneBy)
	iter := iterContaining(endDt, cadence)
	return iter
}

func iterEndDate(iter string) time.Time {
	cadence, valid := iterCadence(iter)
	if !valid {
		panic(fmt.Sprintf("%v was not a valid iteration", iter))
	}
	if cadence == taps.Yearly {
		y, _ := strconv.Atoi(iter)
		return time.Date(y, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	}
	if cadence == taps.Quarterly {
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
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}
	if cadence == taps.Monthly {
		y, _ := strconv.Atoi(iter[0:4])
		mo := strings.ToLower(iter[8:])
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
			panic(fmt.Sprintf("Could not parse to a month: %v", mo))
		}
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}
	if cadence == taps.Biweekly {
		// TODO implement biweekly iterations
		panic("Biweekly iterations not implemented")
	}
	panic(fmt.Sprintf("Could not match the cadence of %v", iter))
}

func iterContaining(date time.Time, cadence taps.Cadence) string {
	y := date.Year()
	if cadence == taps.Yearly {
		return fmt.Sprintf("%v", y)
	} else if cadence == taps.Quarterly {
		if date.Month() <= 3 {
			return fmt.Sprintf("%v Q1", y)
		} else if date.Month() <= 6 {
			return fmt.Sprintf("%v Q2", y)
		} else if date.Month() <= 9 {
			return fmt.Sprintf("%v Q3", y)
		} else {
			return fmt.Sprintf("%v Q4", y)
		}
	} else if cadence == taps.Monthly {
		mos := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		return fmt.Sprintf("%v-%02d %v", y, date.Month(), mos[date.Month()-1])
	}
	panic("Biweekly iterations not implemented")
}

func iterNext(in string) string {
	end := iterEndDate(in)
	start := end.AddDate(0, 0, 1) // Add 1 day to get the first day of the next iteration
	cadence, valid := iterCadence(in)
	if !valid {
		panic(fmt.Sprintf("Could not get cadence of %v: invalid iteration", in))
	}
	return iterContaining(start, cadence)
}

func iterCadence(iter string) (cadence taps.Cadence, valid bool) {
	yrly, _ := regexp.MatchString("^[0-9]{4}$", iter)
	if yrly {
		cadence = taps.Yearly
		valid = true
		return
	}
	qrly, _ := regexp.MatchString("^[0-9]{4} [Qq][1-4]$", iter)
	if qrly {
		cadence = taps.Quarterly
		valid = true
		return
	}
	mnly, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2} [A-Za-z]{3}$", iter)
	if mnly {
		cadence = taps.Monthly
		valid = true
		return
	}
	bwly, _ := regexp.MatchString("^[0-9]{4} [Ww][0-9]{1,2}$", iter)
	if bwly {
		cadence = taps.Biweekly
		valid = true
		return
	}
	valid = false
	return
}
