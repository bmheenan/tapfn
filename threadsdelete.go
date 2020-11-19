package tapfn

import (
	"fmt"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) DeleteThreadHierLinks(anc, child int64) error {
	pas, err := cn.db.GetThreadParentsForAnc(child, anc)
	if err != nil {
		return fmt.Errorf("Could not get parents of %v with ancestor %v: %v", child, anc, err)
	}
	for _, p := range pas {
		des, err := cn.db.GetThreadDes(child)
		if err != nil {
			return fmt.Errorf("Could not get thread descendants: %v", err)
		}
		c := des[child]
		ans, err := cn.db.GetThreadAns(p.ID)
		if err != nil {
			return fmt.Errorf("Could not get thread ancestors: %v", err)
		}
		stks := commonStks(ans, des)
		err = cn.db.DeleteThreadHierLink(p.ID, c.ID)
		if err != nil {
			return fmt.Errorf("Could not unlink %v and %v: %v", p.Name, c.Name, err)
		}
		for _, a := range ans {
			errCTot := cn.recalcCostTot(a)
			if errCTot != nil {
				return fmt.Errorf("Could not update total cost of %v: %v", a.Name, errCTot)
			}
		}
		for stkE := range stks {
			stk, err := cn.db.GetStk(stkE)
			if err != nil {
				return fmt.Errorf("Could not get stakeholder from ans+des stakholders: %v", err)
			}
			iter, err := iterResulting(c.Iter, stk.Cadence)
			if err != nil {
				return fmt.Errorf("Could not get iteration for stakeholder %v: %v", stk.Email, err)
			}
			err = cn.uncrosslinkThreadsForStk(c, p, stk, iter)
			if err != nil {
				return fmt.Errorf("Could not crosslink parent %v with child %v: %v", p.Name, c.Name, err)
			}
			for _, th := range ans {
				if _, ok := th.Stks[stkE]; ok {
					err = cn.recalcCostForStk(th, stk)
					if err != nil {
						return fmt.Errorf("Could not recalc cost of ancestor %v: %v", th.Name, err)
					}
				}
			}
		}
	}
	return nil
}

func commonStks(ans, des map[int64]*taps.Thread) (commonStks map[string](bool)) {
	ancStks := map[string](bool){}
	desStks := map[string](bool){}
	commonStks = map[string](bool){}
	for _, th := range ans {
		for stk := range th.Stks {
			ancStks[stk] = true
		}
	}
	for _, th := range des {
		for stk := range th.Stks {
			desStks[stk] = true
		}
	}
	for stk := range ancStks {
		if _, ok := desStks[stk]; ok {
			commonStks[stk] = true
		}
	}
	return
}

func (cn *cnTapdb) uncrosslinkThreadsForStk(c, p *taps.Thread, stk *taps.Stakeholder, iter string) error {
	chs, errCh := cn.db.GetThreadChildrenByStkIter([]int64{c.ID}, stk.Email, iter)
	if errCh != nil {
		return fmt.Errorf("Could not get child threads: %v", errCh)
	}
	pas, errPa := cn.db.GetThreadParentsByStkIter([]int64{p.ID}, stk.Email, iter)
	if errPa != nil {
		return fmt.Errorf("Could not get parent threads: %v", errPa)
	}
	for _, c := range chs {
		for _, p := range pas {
			errLS := cn.db.DeleteThreadHierLinkForStk(p.ID, c.ID, stk.Email)
			if errLS != nil {
				return fmt.Errorf(
					"Could not link %v with %v for stakeholder %v: %v",
					p.ID,
					c.ID,
					stk.Email,
					errLS,
				)
			}
		}
	}
	return nil
}
