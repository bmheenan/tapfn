package tapfn

import (
	"fmt"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) GetThread(id int64) (th taps.Thread, err error) {
	thp, err := cn.db.GetThread(id)
	if err != nil {
		err = fmt.Errorf("Could not get thread %v: %v", id, err)
		return
	}
	th = *thp
	return
}

func (cn *cnTapdb) GetThreadrowsByStkIter(stk, iter string) (ths []taps.Threadrow, err error) {
	thsp, err := cn.db.GetThreadrowsByStkIter(stk, iter)
	if err != nil {
		err = fmt.Errorf("Could not get threadrows for stakeholder %v and iteration %v: %v", stk, iter, err)
		return
	}
	for _, t := range thsp {
		ths = append(ths, *t)
	}
	return
}

func (cn *cnTapdb) GetThreadrowsByParentIter(parent int64, iter string) (ths []taps.Threadrow, err error) {
	thsp, err := cn.db.GetThreadrowsByParentIter(parent, iter)
	if err != nil {
		err = fmt.Errorf("Could not get threadrows for parent %v and iteration %v: %v", parent, iter, err)
		return
	}
	for _, t := range thsp {
		ths = append(ths, *t)
	}
	return
}
