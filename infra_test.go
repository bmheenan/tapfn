package tapfn

import (
	"fmt"

	"github.com/bmheenan/taps"
)

func setupTest(config string) (
	cn TapController,
	stks map[string]stkInfo,
	iters []string,
	ths map[string]thInfo,
	err error,
) {
	dom := "example.com"
	stks = stkConfigs[config]
	ths = thConfigs[config]
	iters = iterConfigs[config]
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
	for _, si := range stks {
		errStk := cn.NewStk(si.Email, si.Name, si.Abbrev, "#ffffff", "#000000", si.Cadence, si.Parents)
		if errStk != nil {
			err = fmt.Errorf("Could not create stakeholder %v: %v", si.Name, errStk)
			return
		}
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
	Parents []string
	Cost    int
	Owner   string
	Stks    []string
}

var stkConfigs = map[string](map[string]stkInfo){
	"1 stk": {
		"a": stkInfo{
			Email:   "a@example.com",
			Name:    "Person A",
			Abbrev:  "A",
			Cadence: taps.Monthly,
		},
	},
}

var iterConfigs = map[string]([]string){
	"1 stk": {"2020 Oct"},
}

var thConfigs = map[string](map[string]thInfo){
	"1 stk": {},
}
