package tapfn

import (
	"fmt"
	"sort"
	"time"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) GetItersForStk(stkE string) (iters []string, err error) {
	is, errG := cn.db.GetItersForStk(stkE)
	if errG != nil {
		err = fmt.Errorf("Could not get iterations for %v: %v", stkE, errG)
		return
	}
	stk, errStk := cn.db.GetStk(stkE)
	if errStk != nil {
		err = fmt.Errorf("Could not get stakeholder %v: %v", stk, errStk)
		return
	}
	iters, err = cn.itersAddStd(is, stk.Cadence)
	if err != nil {
		err = fmt.Errorf("Could not add standard iterations: %v", err)
		return
	}
	return
}

func (cn *cnTapdb) GetItersForParent(parent int64) (iters []string, err error) {
	is, errG := cn.db.GetItersForParent(parent)
	if errG != nil {
		err = fmt.Errorf("Could not get iterations for %v: %v", parent, errG)
		return
	}
	p, errP := cn.GetThread(parent)
	if errP != nil {
		err = fmt.Errorf("Could not get parent %v: %v", parent, errP)
		return
	}
	iters, err = cn.itersAddStd(is, p.Owner.Cadence)
	if err != nil {
		err = fmt.Errorf("Could not add standard iterations: %v", err)
		return
	}
	return
}

func (cn *cnTapdb) iterCurrent(cadence taps.Cadence) (iter string, err error) {
	var d time.Time
	if cn.timeOverride == (time.Time{}) {
		d = time.Now()
	} else {
		d = cn.timeOverride
	}
	iter, err = iterContaining(d, cadence)
	if err != nil {
		err = fmt.Errorf("Could not get iteration containing the time now: %v", err)
		return
	}
	return
}

func (cn *cnTapdb) itersAddStd(in []string, cadence taps.Cadence) (out []string, err error) {
	current, errCur := cn.iterCurrent(cadence)
	if errCur != nil {
		err = fmt.Errorf("Could not get current iteration: %v", errCur)
		return
	}
	if !strIn(current, in) {
		in = append(in, current)
	}
	next, errNx := iterNext(current)
	if errNx != nil {
		err = fmt.Errorf("Could not get next iteration: %v", errNx)
		return
	}
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
