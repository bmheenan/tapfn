package tapfn

import "fmt"

func (cn *cnTapdb) recalcAllStkCosts(id int64) {
	for _, anc := range cn.db.GetThreadAns(id) {
		for stk := range anc.Stks {
			cn.recalcStkCost(anc.ID, stk)
		}
	}
}

func (cn *cnTapdb) recalcStkCost(id int64, stk string) {
	mbrs, err := cn.db.GetStkDes(stk)
	if err != nil {
		panic(fmt.Sprintf("Could not get stakeholder decendants of %v: %v", stk, err))
	}
	sum := 0
	for _, dec := range cn.db.GetThreadDes(id) {
		if _, ok := mbrs[dec.Owner.Email]; ok {
			sum += dec.CostDir
		}
	}
	err = cn.db.SetCostForStk(id, stk, sum)
	if err != nil {
		panic(fmt.Sprintf("Could not set new cost for %d and %v: %v", id, stk, err))
	}
}
