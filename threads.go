package tapfn

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) NewThread(name, owner, iter string, cost int, parents, children []int64) (int64, error) {
	if name == "" || owner == "" || iter == "" || cost < 0 {
		return 0, errors.New("Name, owner, and iter must be non-blank; cost must be > 0")
	}
	oParts := strings.Split(owner, "@")
	if len(oParts) != 2 {
		return 0, fmt.Errorf("Owner email %v was invalid", owner)
	}
	id, errIn := cn.db.NewThread(name, oParts[1], owner, iter, string(taps.NotStarted), 1, cost)
	if errIn != nil {
		return id, fmt.Errorf("Could not create new thread: %v", errIn)
	}
	errSk := cn.db.NewThreadStkLink(id, owner, oParts[1], iter, math.MaxInt32, true, cost)
	if errSk != nil {
		return id, fmt.Errorf("Could not make owner a stakeholder of the new thread: %v", errSk)
		// TODO: Delete the thread. It's in an invalid state without any stakeholders
	}
	for _, p := range parents {
		errP := cn.NewThreadHierLink(p, id)
		if errP != nil {
			return id, fmt.Errorf("Could not link to parent %v: %v", p, errP)
		}
	}
	for _, c := range children {
		errC := cn.NewThreadHierLink(id, c)
		if errC != nil {
			return id, fmt.Errorf("Could not link to child thread %v: %v", c, errC)
		}
	}
	errB := cn.balanceStk(owner, iter)
	if errB != nil {
		return id, fmt.Errorf("Could not balance threads after creating: %v", errB)
	}
	return id, nil
}

func (cn *cnTapdb) NewThreadHierLink(parent, child int64) error {
	// Don't allow loops of dependencies in threads
	des, errDes := cn.db.GetThreadDes(child)
	if errDes != nil {
		return fmt.Errorf("Could not get thread descendants: %v", errDes)
	}
	ans, errAns := cn.db.GetThreadAns(parent)
	if errAns != nil {
		return fmt.Errorf("Could not get thread ancestors: %v", errAns)
	}
	for a := range ans {
		if _, ok := des[a]; ok {
			return fmt.Errorf("Cannot make %v a parent of %v because that would form a loop", parent, child)
		}
	}
	// Calculate the order and iteration the child will have in the context of this parent thread
	oParts := strings.Split(ans[parent].Owner.Email, "@")
	if len(oParts) != 2 {
		return fmt.Errorf("Parent owner email address is invalid: %v", ans[parent].Owner)
	}
	po, errO := cn.db.GetStk(ans[parent].Owner.Email)
	if errO != nil {
		return fmt.Errorf("Could not get info for parent owner %v: %v", ans[parent].Owner, errO)
	}
	iter, errIt := iterResulting(des[child].Iter, po.Cadence)
	if errIt != nil {
		return fmt.Errorf("Could not convert iteration %v to %v: %v", des[child].Iter, po.Cadence, errIt)
	}
	thOrd, errThO := cn.db.GetOrdBeforeForParent(parent, iter, math.MaxInt32)
	if errThO != nil {
		return fmt.Errorf("Could not get thread order to insert: %v", errThO)
	}
	thOrd = thOrd + ((math.MaxInt32 - thOrd) / 2)
	// Make the main link between the parent and child
	errL := cn.db.NewThreadHierLink(parent, child, ans[parent].Iter, thOrd, oParts[1])
	if errL != nil {
		return fmt.Errorf("Could not link threads: %v", errL)
	}
	errB := cn.balanceParent(parent, iter)
	if errB != nil {
		return fmt.Errorf("Could not balance thread %v for iteration %v after linking: %v", parent, iter, errB)
	}
	// Recalculate the total cost for all ancestor threads, now that they have at least one new descendant
	for _, a := range ans {
		ds, errDs := cn.db.GetThreadDes(a.ID)
		if errDs != nil {
			return fmt.Errorf(
				"Could not get descendants of ancestor %v to calculate its new total cost: %v",
				a.ID,
				errDs,
			)
		}
		rnCost := 0
		for _, d := range ds {
			rnCost += d.CostDir
		}
		errUpC := cn.db.SetCostTot(a.ID, rnCost)
		if errUpC != nil {
			return fmt.Errorf("Could not update total cost for %v: %v", a.ID, errUpC)
		}
	}
	// Calculate changes for each affected stakeholder. Anybody who's a stakeholder of at least one ancestor and at
	// least one descendant will be affected by this link
	ancStks := map[string](bool){}
	desStks := map[string](bool){}
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
	for stkE := range ancStks {
		if _, ok := desStks[stkE]; ok {
			stk, errStk := cn.db.GetStk(stkE)
			if errStk != nil {
				return fmt.Errorf("Could not get stakeholder from ans+des stakholders: %v", errStk)
			}
			stkIter, errISk := iterResulting(des[child].Iter, stk.Cadence)
			if errISk != nil {
				return fmt.Errorf("Could not get iteration for stakeholder %v: %v", stkE, errISk)
			}
			chs, errCh := cn.db.GetThreadChildrenByStkIter([]int64{child}, stkE, stkIter)
			if errCh != nil {
				return fmt.Errorf("Could not get child threads: %v", errCh)
			}
			pas, errPa := cn.db.GetThreadParentsByStkIter([]int64{parent}, stkE, stkIter)
			if errPa != nil {
				return fmt.Errorf("Could not get parent threads: %v", errPa)
			}
			for _, c := range chs {
				for _, p := range pas {
					// Cross link all parents with all children
					errLS := cn.db.NewThreadHierLinkForStk(p.ID, c.ID, stkE, stk.Domain)
					if errLS != nil {
						return fmt.Errorf(
							"Could not link %v with %v for stakeholder %v: %v",
							parent,
							child,
							stkE,
							errLS,
						)
					}
					// Unmark top level for the child; it's now nested under an ancestor thread
					errTop := cn.db.SetTopForStk(c.ID, stkE, false)
					if errTop != nil {
						return fmt.Errorf("Could not mark thread %c as no longer on the top level: %v", c.ID, errTop)
					}
				}
			}
			// Each ancestor thread with this stakeholder needs to have its stakeholder cost recalculated for this
			// iteration
			for _, th := range ans {
				if _, ok := th.Stks[stkE]; ok && th.Stks[stkE].Iter == stkIter {
					ds, errDs := cn.db.GetThreadDes(th.ID)
					if errDs != nil {
						return fmt.Errorf(
							"Could not get descendants of ancestor %v to calculate its new total cost: %v",
							th.ID,
							errDs,
						)
					}
					tms, errTm := cn.db.GetStkDes(stkE)
					if errTm != nil {
						return fmt.Errorf("Could not get team members of %v: %v", stkE, errTm)
					}
					rnCost := 0
					for _, d := range ds {
						if _, ok := tms[d.Owner.Email]; ok {
							dIter, errDI := iterResulting(d.Iter, stk.Cadence)
							if errDI != nil {
								return fmt.Errorf("Could not convert descendant iteration: %v", errDI)
							}
							if dIter == stkIter {
								rnCost += d.CostDir
							}
						}
					}
					errUpC := cn.db.SetCostForStk(th.ID, stkE, rnCost)
					if errUpC != nil {
						return fmt.Errorf(
							"Could not update total cost for thread %v for stakeholder %v: %v",
							th.ID,
							stkE,
							errUpC,
						)
					}
				}
			}
			errBPt := cn.balanceStk(stkE, stkIter)
			if errBPt != nil {
				return fmt.Errorf("Could not balance threads after linking")
			}
		}
	}
	return nil
}
