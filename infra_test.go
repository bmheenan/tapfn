package tapfn

import (
	"fmt"
	"sort"

	"github.com/bmheenan/taps"
)

func setupTest(config string) (
	cn TapController,
	stks map[string]*stkInfo,
	iters []string,
	ths map[string]*thInfo,
	err error,
) {
	dom := "example.com"
	stks = stkCs[config]
	ths = thCs[config]
	iters = iterCs[config]
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
		errStk := cn.NewStk(
			stks[n].Email,
			stks[n].Name,
			stks[n].Abbrev,
			"#ffffff",
			"#000000",
			stks[n].Cadence,
			stks[n].Parents,
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
			Name:    "Person A",
			Abbrev:  "A",
			Cadence: taps.Monthly,
		},
	},
	"1 th": {
		"a": &stkInfo{
			Email:   "a@example.com",
			Name:    "Person A",
			Abbrev:  "A",
			Cadence: taps.Quarterly,
		},
	},
}

var iterCs = map[string]([]string){
	"1 stk": {"2020 Oct"},
	"1 th":  {"2020 Oct"},
}

var thCs = map[string](map[string]*thInfo){
	"1 stk": {},
	"1 th": {
		"A": &thInfo{
			Name:  "A",
			Iter:  iterCs["1 th"][0],
			Cost:  1,
			Owner: stkCs["1 th"]["a"].Email,
		},
	},
}
