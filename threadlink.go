package tapfn

import (
	"errors"
	"fmt"
	"math"
)

func (cn *cnTapdb) ThreadLink(parent, child int64) error {
	if cn.wouldMakeLoop(parent, child) {
		return fmt.Errorf("Cannot link parent %d with child %d: %w", parent, child, ErrWouldMakeLoop)
	}
	p, err := cn.Thread(parent)
	if errors.Is(err, ErrNotFound) {
		return fmt.Errorf("Parent does not exist: %w", err)
	}
	if err != nil {
		panic(err)
	}
	c, err := cn.Thread(child)
	if errors.Is(err, ErrNotFound) {
		return fmt.Errorf("Child does not exist: %w", err)
	}
	if err != nil {
		panic(err)
	}
	iter := iterResulting(c.Iter, p.Owner.Cadence)
	ordLow := cn.db.GetOrdBeforeForParent(parent, iter, math.MaxInt32)
	ord := ((math.MaxInt32 - ordLow) / 2) + ordLow
	cn.db.NewThreadHierLink(parent, child, iter, ord, p.Owner.Domain)
	cn.moveThreadBeforeAns(c.ID)
	cn.recalcAllCostTot(parent)
	cn.balanceParent(parent, iter)
	cn.recalcAllStkCosts(parent)
	cn.recalcPri(p.Owner.Email, p.Stks[p.Owner.Email].Iter)
	return nil
}

func (cn *cnTapdb) ThreadUnlink(parent, child int64) {
	cn.db.DeleteThreadHierLink(parent, child)
	cn.recalcAllCostTot(parent)
	cn.recalcAllStkCosts(parent)
	p, err := cn.Thread(parent)
	if err != nil {
		panic(fmt.Sprintf("Could not get parent: %v", err))
	}
	cn.recalcPri(p.Owner.Email, p.Stks[p.Owner.Email].Iter)
}

func (cn *cnTapdb) wouldMakeLoop(parent, child int64) bool {
	ans := cn.db.GetThreadAns(parent)
	des := cn.db.GetThreadDes(child)
	for a := range ans {
		if _, ok := des[a]; ok {
			return true
		}
	}
	return false
}

func (cn *cnTapdb) moveThreadBeforeAns(thread int64) {
	for _, a := range cn.db.GetThreadAns(thread) {
		if a.ID == thread {
			continue
		}
		th, err := cn.Thread(thread)
		if err != nil {
			panic(fmt.Sprintf("Could not get thread: %v", err))
		}
		for s := range th.Stks {
			if _, ok := a.Stks[s]; !ok {
				continue
			}
			if a.Stks[s].Iter != th.Stks[s].Iter {
				continue
			}
			// Need to refresh thread order; it may have been changed in a previous loop of this thread
			t, err := cn.Thread(thread)
			if err != nil {
				panic(fmt.Sprintf("Could not get thread: %v", err))
			}
			if a.Stks[s].Ord > t.Stks[s].Ord {
				continue
			}
			cn.ThreadMoveForStk(thread, a.ID, s, MoveBeforeRef)
		}
	}

}
