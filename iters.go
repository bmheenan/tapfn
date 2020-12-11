package tapfn

import (
	"fmt"
)

func (cn *cnTapdb) ItersByStk(stkE string) []string {
	is, err := cn.db.GetItersForStk(stkE)
	if err != nil {
		panic(fmt.Sprintf("Could not get iterations for %v: %v", stkE, err))
	}
	stk, err := cn.db.GetStk(stkE)
	if err != nil {
		panic(fmt.Sprintf("Could not get stakeholder %v: %v", stk, err))
	}
	iters, err := cn.itersAddStd(is, stk.Cadence)
	if err != nil {
		panic(fmt.Sprintf("Could not add standard iterations: %v", err))
	}
	return iters
}

func (cn *cnTapdb) ItersByParent(parent int64) []string {
	is, err := cn.db.GetItersForParent(parent)
	if err != nil {
		panic(fmt.Sprintf("Could not get iterations for %v: %v", parent, err))
	}
	p, errP := cn.Thread(parent)
	if errP != nil {
		panic(fmt.Sprintf("Could not get parent %v: %v", parent, err))
	}
	if !strIn(p.Iter, is) {
		is = append(is, p.Iter)
	}
	return is
}

func (cn *cnTapdb) IterOptions(thread int64) []string {
	th, err := cn.db.GetThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread for id %v: %v", thread, err))
	}
	iters := cn.itersBetweenCurrentAnd(th.Iter)
	iters = append(iters, iterNext(iters[len(iters)-1]))
	iters, err = cn.itersAddStd(iters, th.Owner.Cadence)
	if err != nil {
		panic(fmt.Sprintf("Could not add standard iterations: %v", err))
	}
	return iters
}
