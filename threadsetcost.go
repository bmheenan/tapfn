package tapfn

func (cn *cnTapdb) ThreadSetCost(thread int64, cost int) {
	cn.db.SetCostDir(thread, cost)
	for anc := range cn.db.GetThreadAns(thread) {
		cn.recalcCostTot(anc)
	}
	cn.recalcAllStkCosts(thread)
}
