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
			if iterOwnr == d.Iter {
				break
			}
			up := true
			if iterOwnr > d.Iter {
				up = false
			}
			if err != nil {
				return fmt.Errorf("Could not convert %v into owner's cadence for %v: %v", iter, d.Name, err)
			}
			err = cn.db.SetIter(d.ID, iterOwnr)
			if err != nil {
				return fmt.Errorf("Could not set iteration to %v for %v: %v", iterOwnr, th.Name, err)
			}
			for stkE := range d.Stks {
				stk, err := cn.db.GetStk(stkE)
				if err != nil {
					return fmt.Errorf("Could not get stakeholder of descendant thread: %v", err)
				}
				stkIter, err := iterResulting(iter, stk.Cadence)
				if err != nil {
					return fmt.Errorf("Could not convert iteration for stakeholder %v: %v", stkE, err)
				}
				cn.db.SetIterForStk(d.ID, stkE, stkIter)
				if up {
					err = cn.MoveThreadForStk(d.ID, 0, stkE, MoveToEnd)
					if err != nil {
						return fmt.Errorf("Could not move thread to end of iteration: %v", err)
					}
				} else {
					err = cn.MoveThreadForStk(d.ID, 0, stkE, MoveToStart)
					if err != nil {
						return fmt.Errorf("Could not move thread to start of iteration: %v", err)
					}
				}
			}
			for p := range d.Parents {
				pTh, err := cn.db.GetThread(p)
				if err != nil {
					return fmt.Errorf("Could not get parent thread %v: %v", p, err)
				}
				paIter, err := iterResulting(iter, pTh.Owner.Cadence)
				if err != nil {
					return fmt.Errorf("Could not get resulting iteration for parent: %v", err)
				}
				err = cn.db.SetIterForParent(p, d.ID, paIter)
				if err != nil {
					return fmt.Errorf("Could not set iter for parent %v: %v", p, err)
				}
				if up {
					err = cn.MoveThreadForParent(d.ID, 0, p, MoveToEnd)
					if err != nil {
						return fmt.Errorf("Could not move thread to end of iteration: %v", err)
					}
				} else {
					err = cn.MoveThreadForParent(d.ID, 0, p, MoveToStart)
					if err != nil {
						return fmt.Errorf("Could not move thread to start of iteration: %v", err)
					}
				}
			}
			cn.deleteObsoleteStkHierLinks(d)
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
