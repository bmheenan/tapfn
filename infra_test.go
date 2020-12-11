package tapfn

import (
	"fmt"
	"sort"
	"time"

	"github.com/bmheenan/taps"
)

func setupTest(config string) (cn TapController, stks map[string]*stkInfo, ths map[string]*thInfo) {
	dom := "example.com"
	stks = stkCs[config]
	ths = thCs[config]
	cn, err := Init(getTestCredentials())
	if err != nil {
		panic(fmt.Sprintf("Could not get TapController: %v", err))
	}
	cn.DomainClear(dom)
	stkNs := []string{}
	for n := range stks {
		stkNs = append(stkNs, n)
	}
	sort.Strings(stkNs)
	for _, n := range stkNs {
		ps := []string{}
		for _, p := range stks[n].Parents {
			ps = append(ps, stks[p].Email)
		}
		errStk := cn.StkNew(
			stks[n].Email,
			stks[n].Name,
			stks[n].Abbrev,
			"#ffffff",
			"#000000",
			stks[n].Cadence,
			ps,
		)
		if errStk != nil {
			err = fmt.Errorf("Could not create stakeholder %v: %v", stks[n].Name, errStk)
			return
		}
	}
	thNs := []string{}
	for n := range ths {
		thNs = append(thNs, n)
	}
	sort.Strings(thNs)
	for _, n := range thNs {
		ps := []int64{}
		for _, p := range ths[n].Parents {
			ps = append(ps, ths[p].ID)
		}
		id, errTh := cn.ThreadNew(ths[n].Name, ths[n].Owner, ths[n].Iter, ths[n].Cost, ps, nil)
		if errTh != nil {
			err = fmt.Errorf("Could not create thread %v: %v", ths[n].Name, errTh)
			return
		}
		(*ths[n]).ID = id
	}
	if db, ok := cn.(*cnTapdb); ok {
		db.timeOverride = time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)
	}
	return
}

type stkInfo struct {
	Email   string
	Name    string
	Abbrev  string
	Cadence taps.Cadence
	Parents []string
}

type thInfo struct {
	ID      int64
	Name    string
	Iter    string
	Parents []string
	Cost    int
	Owner   string
	Stks    []string
}

var stkCs = map[string](map[string]*stkInfo){
	"blank": {},
	"1 stk": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
			Cadence: taps.Monthly,
		},
	},
	"2 stks w 1 th": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
			Cadence: taps.Monthly,
		},
		"b": &stkInfo{
			Email:   "b@example.com",
			Name:    "Person b",
			Abbrev:  "b",
			Cadence: taps.Monthly,
		},
	},
	"1 th": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
			Cadence: taps.Quarterly,
		},
	},
	"1 stk w ths": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
			Cadence: taps.Monthly,
		},
	},
	"1 stk w ths in diff iters": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
			Cadence: taps.Monthly,
		},
	},
	"s team no ths": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Team a",
			Abbrev:  "a",
			Cadence: taps.Quarterly,
		},
		"aa": &stkInfo{
			Email:   "aa@example.com",
			Name:    "Person aa",
			Abbrev:  "aa",
			Cadence: taps.Monthly,
			Parents: []string{"a"},
		},
		"ab": &stkInfo{
			Email:   "ab@example.com",
			Name:    "Person ab",
			Abbrev:  "ab",
			Cadence: taps.Monthly,
			Parents: []string{"a"},
		},
		"b": &stkInfo{
			Email:   "b@example.com",
			Name:    "Team b",
			Abbrev:  "b",
			Cadence: taps.Quarterly,
		},
		"ba": &stkInfo{
			Email:   "ba@example.com",
			Name:    "Person ba",
			Abbrev:  "ba",
			Cadence: taps.Monthly,
			Parents: []string{"b"},
		},
	},
	"s team": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Team a",
			Abbrev:  "a",
			Cadence: taps.Quarterly,
		},
		"aa": &stkInfo{
			Email:   "aa@example.com",
			Name:    "Person aa",
			Abbrev:  "aa",
			Cadence: taps.Monthly,
			Parents: []string{"a"},
		},
		"ab": &stkInfo{
			Email:   "ab@example.com",
			Name:    "Person ab",
			Abbrev:  "ab",
			Cadence: taps.Monthly,
			Parents: []string{"a"},
		},
		"b": &stkInfo{
			Email:   "b@example.com",
			Name:    "Team b",
			Abbrev:  "b",
			Cadence: taps.Quarterly,
		},
		"ba": &stkInfo{
			Email:   "ba@example.com",
			Name:    "Person ba",
			Abbrev:  "ba",
			Cadence: taps.Monthly,
			Parents: []string{"b"},
		},
	},
	"big tree": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
			Cadence: taps.Monthly,
		},
		"b": &stkInfo{
			Email:   "b@example.com",
			Name:    "Person b",
			Abbrev:  "aa",
			Cadence: taps.Monthly,
		},
	},
}

var thCs = map[string](map[string]*thInfo){
	"blank":         {},
	"1 stk":         {},
	"s team no ths": {},
	"1 th": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020 Q1",
			Cost:  1,
			Owner: stkCs["1 th"]["a"].Email,
		},
	},
	"2 stks w 1 th": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020-12 Dec",
			Cost:  1,
			Owner: stkCs["2 stks w 1 th"]["a"].Email,
		},
	},
	"1 stk w ths": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020-11 Nov",
			Cost:  1,
			Owner: stkCs["1 stk w ths"]["a"].Email,
		},
		"AA": &thInfo{
			Name:    "AA",
			Iter:    "2020-11 Nov",
			Cost:    1,
			Owner:   stkCs["1 stk w ths"]["a"].Email,
			Parents: []string{"A"},
		},
		"AB": &thInfo{
			Name:    "AB",
			Iter:    "2020-11 Nov",
			Cost:    1,
			Owner:   stkCs["1 stk w ths"]["a"].Email,
			Parents: []string{"A"},
		},
		"AC": &thInfo{
			Name:    "AC",
			Iter:    "2020-11 Nov",
			Cost:    1,
			Owner:   stkCs["1 stk w ths"]["a"].Email,
			Parents: []string{"A"},
		},
	},
	"1 stk w ths in diff iters": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020-10 Oct",
			Cost:  1,
			Owner: stkCs["1 stk w ths in diff iters"]["a"].Email,
		},
		"B": &thInfo{
			Name:  "B",
			Iter:  "2020-11 Nov",
			Cost:  1,
			Owner: stkCs["1 stk w ths in diff iters"]["a"].Email,
		},
	},
	"s team": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020 Q4",
			Cost:  0,
			Owner: stkCs["s team"]["a"].Email,
		},
		"AA": &thInfo{
			Name:    "AA",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["aa"].Email,
			Parents: []string{"A"},
		},
		"AB": &thInfo{
			Name:    "AB",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["ab"].Email,
			Parents: []string{"A"},
		},
		"AC": &thInfo{
			Name:    "AC",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["ab"].Email,
			Parents: []string{"A"},
		},
		"B": &thInfo{
			Name:  "B",
			Iter:  "2020-12 Dec",
			Cost:  10,
			Owner: stkCs["s team"]["aa"].Email,
		},
	},
	"big tree": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020-10 Oct",
			Cost:  0,
			Owner: stkCs["s team"]["a"].Email,
		},
		"AA": &thInfo{
			Name:    "AA",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["b"].Email,
			Parents: []string{"A"},
		},
		"AAA": &thInfo{
			Name:    "AAA",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["a"].Email,
			Parents: []string{"AA"},
		},
		"AB": &thInfo{
			Name:    "AB",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["a"].Email,
			Parents: []string{"A"},
		},
		"AC": &thInfo{
			Name:    "AC",
			Iter:    "2020-10 Oct",
			Cost:    5,
			Owner:   stkCs["s team"]["b"].Email,
			Parents: []string{"A"},
		},
		"B": &thInfo{
			Name:  "B",
			Iter:  "2020-10 Oct",
			Cost:  10,
			Owner: stkCs["s team"]["a"].Email,
		},
		"C": &thInfo{
			Name:  "C",
			Iter:  "2020-10 Oct",
			Cost:  10,
			Owner: stkCs["s team"]["b"].Email,
		},
	},
}
