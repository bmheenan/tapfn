package tapfn

import (
	"fmt"
	"math"
)

// MoveTo specifies the different anchors you can move a thread to within an iteration
type MoveTo int

const (
	// MoveToStart moves the thread to the beginning of the iteration, igoring the reference
	MoveToStart = iota
	// MoveToEnd moves the thread to the end of the iteration, ignoring the reference
	MoveToEnd
	// MoveBeforeRef moves the thread to right before the given reference
	MoveBeforeRef
)

func (cn *cnTapdb) MoveThreadForParent(thread, reference, parent int64, moveTo MoveTo) error {
	var nOrd int
	t, err := cn.db.GetThread(thread)
	if err != nil {
		return fmt.Errorf("Could not get thread for thread %v: %v", thread, err)
	}
	p, err := cn.db.GetThread(parent)
	if err != nil {
		return fmt.Errorf("Could not get thread for parent thread %v: %v", parent, err)
	}
	po, err := cn.db.GetStk(p.Owner.Email)
	if err != nil {
		return fmt.Errorf("Could not get stakeholder %v from parent owner: %v", p.Owner, err)
	}
	ti, err := iterResulting(t.Iter, po.Cadence)
	if err != nil {
		return fmt.Errorf("Could not get thread's iteration by the parent owner's cadence: %v", err)
	}
	switch moveTo {
	case MoveToStart:
		ordA, err := cn.db.GetOrdAfterForParent(parent, ti, 0)
		if err != nil {
			return fmt.Errorf("Could not get order of first thread in this iteration under this parent: %v", err)
		}
		nOrd = ordA / 2
	case MoveToEnd:
		ordB, err := cn.db.GetOrdBeforeForParent(parent, ti, math.MaxInt32)
		if err != nil {
			return fmt.Errorf("Could not get order of last thread in iteration under this parent: %v", err)
		}
		nOrd = ordB + ((math.MaxInt32 - ordB) / 2)
	case MoveBeforeRef:
		r, err := cn.db.GetThread(reference)
		if err != nil {
			return fmt.Errorf("Could not get threadrel for reference thread %v: %v", reference, err)
		}
		ri, err := iterResulting(r.Iter, po.Cadence)
		if err != nil {
			return fmt.Errorf("Could not get reference's iteration by the parent owner's cadence: %v", err)
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
		ordB, err := cn.db.GetOrdBeforeForParent(parent, ti, r.Parents[parent].Ord)
		if err != nil {
			return fmt.Errorf("Could not get order of thread before reference under this parent: %v", err)
		}
		nOrd = ordB + ((r.Parents[parent].Ord - ordB) / 2)
	}
	err = cn.db.SetOrdForParent(thread, parent, nOrd)
	if err != nil {
		return fmt.Errorf("Could not set new order for thread: %v", err)
	}
	return nil
}

func (cn *cnTapdb) MoveThreadForStk(thread, reference int64, stkE string, moveTo MoveTo) error {
	var nOrd int
	t, err := cn.db.GetThread(thread)
	if err != nil {
		return fmt.Errorf("Could not get thread %v: %v", thread, err)
	}
	stk, err := cn.db.GetStk(stkE)
	if err != nil {
		return fmt.Errorf("Could not get stakeholder %v: %v", stk, err)
	}
	ti, err := iterResulting(t.Iter, stk.Cadence)
	if err != nil {
		return fmt.Errorf("Could not get thread's iteration by the stakeholder's cadence: %v", err)
	}
	switch moveTo {
	case MoveToStart:
		ordA, err := cn.db.GetOrdAfterForStk(stkE, ti, 0)
		if err != nil {
			return fmt.Errorf("Could not get order of first thread in the iteration for this stakeholder: %v", err)
		}
		nOrd = ordA / 2
	case MoveToEnd:
		ordB, err := cn.db.GetOrdBeforeForStk(stkE, ti, math.MaxInt32)
		if err != nil {
			return fmt.Errorf("Could not get order of last thread in iteration for this stakeholder: %v", err)
		}
		nOrd = ordB + ((math.MaxInt32 - ordB) / 2)
	case MoveBeforeRef:
		r, err := cn.db.GetThread(reference)
		if err != nil {
			return fmt.Errorf("Could not get threadrel for reference thread %v: %v", reference, err)
		}
		ri, err := iterResulting(r.Iter, stk.Cadence)
		if err != nil {
			return fmt.Errorf("Could not get reference's iteration by the stakeholder's cadence: %v", err)
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
		ordB, err := cn.db.GetOrdBeforeForStk(stkE, ti, r.Stks[stkE].Ord)
		if err != nil {
			return fmt.Errorf("Could not get order of thread before reference under this parent: %v", err)
		}
		nOrd = ordB + ((r.Stks[stkE].Ord - ordB) / 2)
	}
	err = cn.db.SetOrdForStk(thread, stkE, nOrd)
	if err != nil {
		return fmt.Errorf("Could not set new order for thread: %v", err)
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
