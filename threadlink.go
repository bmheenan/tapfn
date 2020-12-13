package tapfn

import (
	"errors"
	"fmt"
	"math"

	"github.com/bmheenan/taps"
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
	cn.moveThreadBeforeAns(c)
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

func (cn *cnTapdb) moveThreadBeforeAns(thread taps.Thread) {
	for _, a := range cn.db.GetThreadAns(thread.ID) {
		if a.ID == thread.ID {
			continue
		}
		for s := range thread.Stks {
			if _, ok := a.Stks[s]; !ok {
				continue
			}
			if a.Stks[s].Iter != thread.Stks[s].Iter {
				continue
			}
			if a.Stks[s].Ord > thread.Stks[s].Ord {
				continue
			}
			cn.ThreadMoveForStk(thread.ID, a.ID, s, MoveBeforeRef)
		}
	}

}
