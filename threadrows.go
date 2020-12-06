package tapfn

import "github.com/bmheenan/taps"

func (cn *cnTapdb) ThreadrowsByStkIter(stk, iter string) []taps.Threadrow {
	return cn.db.GetThreadrowsByStkIter(stk, iter)
}

func (cn *cnTapdb) ThreadrowsByParentIter(parent int64, iter string) []taps.Threadrow {
	return cn.db.GetThreadrowsByParentIter(parent, iter)
}

func (cn *cnTapdb) ThreadrowsByChild(child int64) []taps.Threadrow {
	return cn.db.GetThreadrowsByChild(child)
}
