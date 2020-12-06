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
	for a := range cn.db.GetThreadAns(parent) {
		cn.recalcCostTot(a)
	}
	cn.balanceParent(parent, iter)
	cn.recalcAllStkCosts(parent)
	return nil
}

func (cn *cnTapdb) ThreadUnlink(parent, child int64) {
	cn.db.DeleteThreadHierLink(parent, child)
	for a := range cn.db.GetThreadAns(parent) {
		cn.recalcCostTot(a)
	}
	cn.recalcAllStkCosts(parent)
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
