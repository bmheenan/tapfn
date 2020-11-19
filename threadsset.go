package tapfn

import (
	"fmt"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) SetThreadIter(thID int64, iter string) error {
	th, err := cn.db.GetThread(thID)
	if err != nil {
		return fmt.Errorf("Could not get thread %v: %v", thID, err)
	}
	des, err := cn.db.GetThreadDes(th.ID)
	if err != nil {
		return fmt.Errorf("Could not get thread descendants: %v", err)
	}
	oldIter := th.Iter
	affectedAns := map[int64](*taps.Thread){}
	for _, d := range des {
		if d.Iter == oldIter {
			iterOwnr, err := iterRequired(iter, d.Owner.Cadence)
			if err != nil {
				return fmt.Errorf("Could not convert %v into owner's cadence for %v: %v", iter, d.Name, err)
			}
			err = cn.db.SetIter(d.ID, iterOwnr)
			if err != nil {
				return fmt.Errorf("Could not set iteration to %v for %v: %v", iterOwnr, th.Name, err)
			}
			ans, errA := cn.db.GetThreadAns(d.ID)
			if errA != nil {
				return fmt.Errorf("Could not get ancestors of %v: %v", d.Name, err)
			}
			for _, a := range ans {
				affectedAns[a.ID] = a
			}
		}
	}
	for _, a := range affectedAns {
		for stkE := range a.Stks {
			stk, err := cn.db.GetStk(stkE)
			if err != nil {
				return fmt.Errorf("Could not get stakeholder of %v: %v", a.Name, err)
			}
			cn.recalcCostForStk(a, stk)
		}
	}
	return nil
}
