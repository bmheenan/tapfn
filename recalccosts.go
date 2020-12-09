package tapfn

import (
	"fmt"

	"github.com/bmheenan/taps"
)

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
	s, err := cn.Stk(stk)
	if err != nil {
		panic(fmt.Sprintf("Could not get stakeholder: %v", err))
	}
	th, err := cn.Thread(id)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread: %v", err))
	}
	sum := 0
	for _, d := range cn.db.GetThreadDes(id) {
		if _, ok := mbrs[d.Owner.Email]; !ok {
			continue
		}
		if d.State != taps.NotStarted && d.State != taps.InProgress {
			continue
		}
		if iterResulting(d.Iter, s.Cadence) != iterResulting(th.Iter, s.Cadence) {
			continue
		}
		sum += d.CostDir
	}
	err = cn.db.SetCostForStk(id, stk, sum)
	if err != nil {
		panic(fmt.Sprintf("Could not set new cost for %d and %v: %v", id, stk, err))
	}
}

func (cn *cnTapdb) recalcAllCostTot(id int64) {
	for a := range cn.db.GetThreadAns(id) {
		cn.recalcCostTot(a)
	}
}

func (cn *cnTapdb) recalcCostTot(id int64) {
	ds := cn.db.GetThreadDes(id)
	rnCost := 0
	for _, d := range ds {
		if d.State != taps.NotStarted && d.State != taps.InProgress {
			continue
		}
		rnCost += d.CostDir
	}
	cn.db.SetCostTot(id, rnCost)
}
