package tapfn

func (cn *cnTapdb) ThreadSetCost(thread int64, cost int) {
	cn.db.SetCostDir(thread, cost)
	cn.recalcAllCostTot(thread)
	cn.recalcAllStkCosts(thread)
}
