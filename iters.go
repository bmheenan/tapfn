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
