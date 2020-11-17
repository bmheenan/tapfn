package tapfn

import (
	"fmt"
	"sort"
	"time"

	"github.com/bmheenan/taps"
)

func setupTest(config string) (
	cn TapController,
	stks map[string]*stkInfo,
	ths map[string]*thInfo,
	err error,
) {
	dom := "example.com"
	stks = stkCs[config]
	ths = thCs[config]
	var errCn error
	cn, errCn = Init(getTestCredentials())
	if errCn != nil {
		err = fmt.Errorf("Could not get TapController: %v", errCn)
		return
	}
	errClr := cn.ClearDomain(dom)
	if errClr != nil {
		err = fmt.Errorf("Could not clear domain: %v", errClr)
		return
	}
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
		errStk := cn.NewStk(
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
		id, errTh := cn.NewThread(ths[n].Name, ths[n].Owner, ths[n].Iter, ths[n].Cost, ps, nil)
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
	"1 stk": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person a",
			Abbrev:  "a",
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
}

var thCs = map[string](map[string]*thInfo){
	"1 stk": {},
	"1 th": {
		"A": &thInfo{
			Name:  "A",
			Iter:  "2020 Q1",
			Cost:  1,
			Owner: stkCs["1 th"]["a"].Email,
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
}
