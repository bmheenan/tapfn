package tapfn

import (
	"fmt"
)

func (cn *cnTapdb) ThreadSetCost(thread int64, cost int) {
	cn.db.SetCostDir(thread, cost)
	cn.recalcAllCostTot(thread)
	cn.recalcAllStkCosts(thread)
	th, err := cn.Thread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread: %v", err))
	}
	cn.recalcPri(th.Owner.Email, th.Stks[th.Owner.Email].Iter)
}
