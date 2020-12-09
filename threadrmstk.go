package tapfn

func (cn *cnTapdb) ThreadRemoveStk(thread int64, stk string) {
	cn.db.DeleteThreadStkLink(thread, stk)
	cn.recalcAllStkCosts(thread)
}
