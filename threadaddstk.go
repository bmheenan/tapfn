package tapfn

import "math"

func (cn *cnTapdb) ThreadAddStk(id int64, stk string) {
	s, err := cn.Stk(stk)
	if err != nil {
		return
	}
	th, err := cn.Thread(id)
	if err != nil {
		return
	}
	iter := iterResulting(th.Iter, s.Cadence)
	ordLow := cn.db.GetOrdBeforeForStk(stk, iter, math.MaxInt32)
	ord := ((math.MaxInt32 - ordLow) / 2) + ordLow
	cn.db.NewThreadStkLink(id, stk, s.Domain, iter, ord, 0)
	cn.balanceStk(stk, iter)
	cn.recalcAllStkCosts(id)
}
