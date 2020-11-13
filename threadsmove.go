package tapfn

import (
	"fmt"
	"math"
)

// MoveThreadParent moves the thread with id `thread` immediately before the thread with id `reference`, in the context
// of the thread with id `parent`. `parent` must be a parent of the other two threads, and the two threads must be in
// the same iteration
// If `reference` = 0, `thread` will be moved to the end of the iteration
func (cn *cnTapdb) MoveThreadForParent(thread, reference, parent int64) error {
	var nOrd int
	t, errT := cn.db.GetThread(thread)
	if errT != nil {
		return fmt.Errorf("Could not get thread for thread %v: %v", thread, errT)
	}
	p, errP := cn.db.GetThread(parent)
	if errP != nil {
		return fmt.Errorf("Could not get thread for parent thread %v: %v", parent, errP)
	}
	po, errPO := cn.db.GetStk(p.Owner.Email)
	if errPO != nil {
		return fmt.Errorf("Could not get stakeholder %v from parent owner: %v", p.Owner, errPO)
	}
	ti, errTI := iterResulting(t.Iter, po.Cadence)
	if errTI != nil {
		return fmt.Errorf("Could not get thread's iteration by the parent owner's cadence: %v", errTI)
	}
	if reference == 0 {
		// Put thread at the end of the iteration
		ordB, errOrdB := cn.db.GetOrdBeforeForParent(parent, ti, math.MaxInt32)
		if errOrdB != nil {
			return fmt.Errorf("Could not get order of last thread in iteration under this parent: %v", errOrdB)
		}
		nOrd = ordB + ((math.MaxInt32 - ordB) / 2)
	} else {
		// Put thread immediately before reference, if they're in the same iteration
		r, errR := cn.db.GetThread(reference)
		if errR != nil {
			return fmt.Errorf("Could not get threadrel for reference thread %v: %v", reference, errR)
		}
		ri, errRI := iterResulting(r.Iter, po.Cadence)
		if errRI != nil {
			return fmt.Errorf("Could not get reference's iteration by the parent owner's cadence: %v", errRI)
		}
		if ti != ri {
			return fmt.Errorf(
				"Cannot move thread %v (iteration %v=%v) before thread %v (iteration %v:=%v): different iterations",
				thread,
				t.Iter,
				ti,
				reference,
				r.Iter,
				ri,
			)
		}
		ordB, errOrdB := cn.db.GetOrdBeforeForParent(parent, ti, r.Parents[parent].Ord)
		if errOrdB != nil {
			return fmt.Errorf("Could not get order of thread before reference under this parent: %v", errOrdB)
		}
		nOrd = ordB + ((r.Parents[parent].Ord - ordB) / 2)
	}
	errM := cn.db.SetOrdForParent(thread, parent, nOrd)
	if errM != nil {
		return fmt.Errorf("Could not set new order for thread: %v", errM)
	}
	return nil
}

// MoveThreadStakeholder moves the thread with id `thread` immediately before the thread with id `reference`, as long as
// `stakeholder` is a stakeholder of both, and they appear in the same iteration for that stakeholder
// If `reference` = 0, `thread` will be moved to the end of the iteration
func (cn *cnTapdb) MoveThreadForStk(thread, reference int64, stkE string) error {
	var nOrd int
	t, errT := cn.db.GetThread(thread)
	if errT != nil {
		return fmt.Errorf("Could not get thread %v: %v", thread, errT)
	}
	stk, errPT := cn.db.GetStk(stkE)
	if errPT != nil {
		return fmt.Errorf("Could not get stakeholder %v: %v", stk, errPT)
	}
	ti, errTI := iterResulting(t.Iter, stk.Cadence)
	if errTI != nil {
		return fmt.Errorf("Could not get thread's iteration by the stakeholder's cadence: %v", errTI)
	}
	if reference == 0 {
		// Put thread at the end of the iteration
		ordB, errOrdB := cn.db.GetOrdBeforeForStk(stkE, ti, math.MaxInt32)
		if errOrdB != nil {
			return fmt.Errorf("Could not get order of last thread in iteration for this stakeholder: %v", errOrdB)
		}
		nOrd = ordB + ((math.MaxInt32 - ordB) / 2)
	} else {
		// Put thread immediately before reference, if they're in the same iteration
		r, errR := cn.db.GetThread(reference)
		if errR != nil {
			return fmt.Errorf("Could not get threadrel for reference thread %v: %v", reference, errR)
		}
		ri, errRI := iterResulting(r.Iter, stk.Cadence)
		if errRI != nil {
			return fmt.Errorf("Could not get reference's iteration by the stakeholder's cadence: %v", errRI)
		}
		if ti != ri {
			return fmt.Errorf(
				"Cannot move thread %v (iteration %v=%v) before thread %v (iteration %v:=%v): different iterations",
				thread,
				t.Iter,
				ti,
				reference,
				r.Iter,
				ri,
			)
		}
		ordB, errOrdB := cn.db.GetOrdBeforeForStk(stkE, ti, r.Stks[stkE].Ord)
		if errOrdB != nil {
			return fmt.Errorf("Could not get order of thread before reference under this parent: %v", errOrdB)
		}
		nOrd = ordB + ((r.Stks[stkE].Ord - ordB) / 2)
	}
	errM := cn.db.SetOrdForStk(thread, stkE, nOrd)
	if errM != nil {
		return fmt.Errorf("Could not set new order for thread: %v", errM)
	}
	return nil
}

func (cn *cnTapdb) balanceParent(parent int64, iter string) error {
	ths, errThs := cn.db.GetThreadsByParentIter(parent, iter)
	if errThs != nil {
		return fmt.Errorf("Could not get threads under thread %v in iteration %v: %v", parent, iter, errThs)
	}
	step := math.MaxInt32 / (len(ths) + 1)
	for i, th := range ths {
		errO := cn.db.SetOrdForParent(th.ID, parent, step*(i+1))
		if errO != nil {
			return fmt.Errorf("Could not set thread order: %v", errO)
		}
	}
	return nil
}

func (cn *cnTapdb) balanceStk(stk, iter string) error {
	ths, errThs := cn.db.GetThreadsByStkIter(stk, iter)
	if errThs != nil {
		return fmt.Errorf("Could not get threads for %v in iteration %v: %v", stk, iter, errThs)
	}
	step := math.MaxInt32 / (len(ths) + 1)
	for i, th := range ths {
		errO := cn.db.SetOrdForStk(th.ID, stk, step*(i+1))
		if errO != nil {
			return fmt.Errorf("Could not set thread order: %v", errO)
		}
	}
	return nil
}
