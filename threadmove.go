package tapfn

import (
	"fmt"
	"math"
)

func (cn *cnTapdb) ThreadMoveForStk(thread, reference int64, stkE string, moveTo MoveTo) {
	var nOrd int
	t, err := cn.db.GetThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread %v: %v", thread, err))
	}
	stk, err := cn.db.GetStk(stkE)
	if err != nil {
		panic(fmt.Sprintf("Could not get stakeholder %v: %v", stk, err))
	}
	ti := iterResulting(t.Iter, stk.Cadence)
	switch moveTo {
	case MoveToStart:
		ordA := cn.db.GetOrdAfterForStk(stkE, ti, 0)
		nOrd = ordA / 2
	case MoveToEnd:
		ordB := cn.db.GetOrdBeforeForStk(stkE, ti, math.MaxInt32)
		nOrd = ordB + ((math.MaxInt32 - ordB) / 2)
	case MoveBeforeRef:
		var lowestID int64
		lowestOrd := math.MaxInt32
		for _, dec := range cn.db.GetThreadDes(reference) {
			if dec.Stks[stkE].Ord < lowestOrd {
				lowestOrd = dec.Stks[stkE].Ord
				lowestID = dec.ID
			}
		}
		r, err := cn.db.GetThread(lowestID)
		if err != nil {
			panic(fmt.Sprintf("Could not get threadrel for dec reference thread %v: %v", lowestID, err))
		}
		ri := iterResulting(r.Iter, stk.Cadence)
		if ti != ri {
			panic(fmt.Sprintf("Can't move thread %v (%v=%v) before %v (%v:=%v): different iterations", thread, t.Iter, ti, lowestID, r.Iter, ri))
		}
		ordB := cn.db.GetOrdBeforeForStk(stkE, ti, r.Stks[stkE].Ord)
		nOrd = ordB + ((r.Stks[stkE].Ord - ordB) / 2)
	}
	err = cn.db.SetOrdForStk(thread, stkE, nOrd)
	if err != nil {
		panic(fmt.Sprintf("Could not set new order for thread: %v", err))
	}
	if t.Owner.Email == stkE {
		cn.recalcPri(stkE, t.Stks[stkE].Iter)
	}
}

func (cn *cnTapdb) ThreadMoveForParent(thread, reference, parent int64, moveTo MoveTo) {
	var nOrd int
	t, err := cn.db.GetThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread for thread %v: %v", thread, err))
	}
	p, err := cn.db.GetThread(parent)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread for parent thread %v: %v", parent, err))
	}
	po, err := cn.db.GetStk(p.Owner.Email)
	if err != nil {
		panic(fmt.Sprintf("Could not get stakeholder %v from parent owner: %v", p.Owner, err))
	}
	ti := iterResulting(t.Iter, po.Cadence)
	switch moveTo {
	case MoveToStart:
		ordA := cn.db.GetOrdAfterForParent(parent, ti, 0)
		nOrd = ordA / 2
	case MoveToEnd:
		ordB := cn.db.GetOrdBeforeForParent(parent, ti, math.MaxInt32)
		nOrd = ordB + ((math.MaxInt32 - ordB) / 2)
	case MoveBeforeRef:
		r, err := cn.db.GetThread(reference)
		if err != nil {
			panic(fmt.Sprintf("Could not get threadrel for reference thread %v: %v", reference, err))
		}
		ri := iterResulting(r.Iter, po.Cadence)
		if ti != ri {
			panic(fmt.Sprintf(
				"Cannot move thread %v (iteration %v=%v) before thread %v (iteration %v:=%v): different iterations",
				thread,
				t.Iter,
				ti,
				reference,
				r.Iter,
				ri,
			))
		}
		ordB := cn.db.GetOrdBeforeForParent(parent, ti, r.Parents[parent].Ord)
		nOrd = ordB + ((r.Parents[parent].Ord - ordB) / 2)
	}
	cn.db.SetOrdForParent(thread, parent, nOrd)
}
