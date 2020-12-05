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
	des, err := cn.db.GetThreadDes(child)
	if err != nil {
		return fmt.Errorf("Could not get thread descendants: %v", err)
	}
	c := des[child]
	ans, err := cn.db.GetThreadAns(parent)
	if err != nil {
		return fmt.Errorf("Could not get thread ancestors: %v", err)
	}
	p := ans[parent]
	if _, ok := c.Parents[p.ID]; ok {
		return nil
	}
	if cn.wouldMakeLoop(ans, des) {
		return fmt.Errorf("Cannot make %v a parent of %v because that would make a loop", parent, child)
	}
	err = cn.newThreadHierLinkForParent(p, c)
	if err != nil {
		return fmt.Errorf("Could not link %v to %v: %v", parent, child, err)
	}
	for _, a := range ans {
		err = cn.recalcCostTot(a)
		if err != nil {
			return fmt.Errorf("Could not update total cost of %v: %v", a.Name, err)
		}
	}
	err = cn.makeThreadStkHierLinksForPair(p, c, ans, des)
	if err != nil {
		return fmt.Errorf("Could not make stakeholder thread hierarchy links: %v", err)
	}
	return nil
}

// Returns if linking a thread in `ans` to a thread in `des` would create a loop
func (cn *cnTapdb) wouldMakeLoop(ans, des map[int64]*taps.Thread) bool {
	for a := range ans {
		if _, ok := des[a]; ok {
			return true
		}
	}
	return false
}

// Forms the core parent/child link, but does not handle stakeholders' contexts
func (cn *cnTapdb) newThreadHierLinkForParent(parent, child *taps.Thread) error {
	oParts := strings.Split(parent.Owner.Email, "@")
	if len(oParts) != 2 {
		return fmt.Errorf("Parent owner email address is invalid: %v", parent.Owner)
	}
	iter, errIt := iterResulting(child.Iter, parent.Owner.Cadence)
	if errIt != nil {
		return fmt.Errorf(
			"Could not convert iteration %v to %v: %v",
			child.Iter,
			parent.Owner.Cadence,
			errIt,
		)
	}
	ord, errOrd := cn.db.GetOrdBeforeForParent(parent.ID, iter, math.MaxInt32)
	if errOrd != nil {
		return fmt.Errorf("Could not get thread order to insert: %v", errOrd)
	}
	ord = ord + ((math.MaxInt32 - ord) / 2)
	errL := cn.db.NewThreadHierLink(parent.ID, child.ID, iter, ord, oParts[1])
	if errL != nil {
		return fmt.Errorf("Could not link threads: %v", errL)
	}
	errB := cn.balanceParent(parent.ID, iter)
	if errB != nil {
		return fmt.Errorf("Could not balance thread %v for iteration %v after linking: %v", parent, iter, errB)
	}
	return nil
}

// Removes all stk thread hier links between a thread and its ancestors/descedants if those threads are in different
// iterations or at least one of them doesn't have the given stakeholder
func (cn *cnTapdb) deleteObsoleteStkHierLinks(th *taps.Thread) error {
	for stkE := range th.Stks {
		ans, err := cn.db.GetThreadAns(th.ID)
		if err != nil {
			return fmt.Errorf("Could not get thread ancestors: %v", err)
		}
		des, err := cn.db.GetThreadDes(th.ID)
		if err != nil {
			return fmt.Errorf("Could not get thread descendants: %v", err)
		}
		type rel struct {
			thread *taps.Thread
			anc    bool
		}
		rels := []rel{}
		for _, a := range ans {
			rels = append(rels, rel{
				thread: a,
				anc:    true,
			})
		}
		for _, d := range des {
			rels = append(rels, rel{
				thread: d,
				anc:    false,
			})
		}
		for _, r := range rels {
			del := false
			if _, ok := r.thread.Stks[stkE]; !ok {
				del = true
			}
			if _, ok := th.Stks[stkE]; !ok {
				del = true
			}
			if !del && th.Stks[stkE].Iter != r.thread.Stks[stkE].Iter {
				del = true
			}
			if del && r.anc {
				err = cn.db.DeleteThreadHierLinkForStk(r.thread.ID, th.ID, stkE)
				if err != nil {
					return fmt.Errorf("Could not delete thread hier link: %v", err)
				}
			} else if del {
				err = cn.db.DeleteThreadHierLinkForStk(th.ID, r.thread.ID, stkE)
				if err != nil {
					return fmt.Errorf("Could not delete thread hier link: %v", err)
				}
			}
		}
	}

	return nil
}

func (cn *cnTapdb) makeAllThreadStkHierLinks(th *taps.Thread) error {
	return errors.New("Not implemented")
}

func (cn *cnTapdb) makeThreadStkHierLinksForPair(
	p *taps.Thread,
	c *taps.Thread,
	ans map[int64]*taps.Thread,
	des map[int64]*taps.Thread,
) error {
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
			err := cn.makeThreadHierLinksForStk(p, c, stkE)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cn *cnTapdb) makeThreadHierLinksForStk(p *taps.Thread, c *taps.Thread, stkE string) error {
	stk, err := cn.db.GetStk(stkE)
	if err != nil {
		return fmt.Errorf("Could not get stakeholder from ans+des stakholders: %v", err)
	}
	iter, err := iterResulting(c.Iter, stk.Cadence)
	if err != nil {
		return fmt.Errorf("Could not get iteration for stakeholder %v: %v", stk.Email, err)
	}
	ans, err := cn.db.GetThreadAns(p.ID)
	if err != nil {
		return fmt.Errorf("Could not get ancestors of thread: %v", err)
	}
	err = cn.crosslinkThreadsForStk(c, p, stk, iter)
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
	err = cn.balanceStk(stk.Email, iter)
	if err != nil {
		return fmt.Errorf("Could not balance threads after linking: %v", err)
	}
	return nil
}

// recalcCostTot updates the total cost of `thread` to be the sum of its direct cost and all its descendants
func (cn *cnTapdb) recalcCostTot(thread *taps.Thread) error {
	ds, errDs := cn.db.GetThreadDes(thread.ID)
	if errDs != nil {
		return fmt.Errorf(
			"Could not get descendants of %v to calculate its new total cost: %v",
			thread.ID,
			errDs,
		)
	}
	rnCost := 0
	for _, d := range ds {
		rnCost += d.CostDir
	}
	errUpC := cn.db.SetCostTot(thread.ID, rnCost)
	if errUpC != nil {
		return fmt.Errorf(
			"Could not update total cost for thread %v: %v",
			thread.ID,
			errUpC,
		)
	}
	return nil
}

// Links `c` and `p` together for stakeholder `stk` in iteration `iter`. Will form multiple links where necessary
func (cn *cnTapdb) crosslinkThreadsForStk(c, p *taps.Thread, stk *taps.Stakeholder, iter string) error {
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
			errLS := cn.db.NewThreadHierLinkForStk(p.ID, c.ID, stk.Email, stk.Domain)
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

// recalcCostForStkIter updates the stakeholder cost for the given `thread`, `stk` - the cost
// that that stakeholder (and all team members) own within the iteration of the thread
func (cn *cnTapdb) recalcCostForStk(thread *taps.Thread, stk *taps.Stakeholder) error {
	ds, errDs := cn.db.GetThreadDes(thread.ID)
	if errDs != nil {
		return fmt.Errorf(
			"Could not get descendants of ancestor %v to calculate its new total cost: %v",
			thread.ID,
			errDs,
		)
	}
	tms, errTm := cn.db.GetStkDes(stk.Email)
	if errTm != nil {
		return fmt.Errorf("Could not get team members of %v: %v", stk.Email, errTm)
	}
	rnCost := 0
	for _, d := range ds {
		if _, ok := tms[d.Owner.Email]; ok {
			dIter, errDI := iterResulting(d.Iter, stk.Cadence)
			if errDI != nil {
				return fmt.Errorf("Could not convert descendant iteration: %v", errDI)
			}
			if dIter == thread.Stks[stk.Email].Iter {
				rnCost += d.CostDir
			}
		}
	}
	errUpC := cn.db.SetCostForStk(thread.ID, stk.Email, rnCost)
	if errUpC != nil {
		return fmt.Errorf(
			"Could not update total cost for thread %v for stakeholder %v: %v",
			thread.ID,
			stk.Email,
			errUpC,
		)
	}
	return nil
}
