package tapfn

import (
	"fmt"
	"strings"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) ThreadSetState(thread int64, state taps.State) {
	cn.db.SetState(thread, state)
	thsAffected := map[int64]bool{}
	stkItersAffected := map[string]bool{}
	th, err := cn.Thread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread: %v", err))
	}
	switch state {
	case taps.NotStarted:
		for _, a := range cn.db.GetThreadAns(thread) {
			if a.State == taps.Done || a.State == taps.Closed || a.State == taps.Archived {
				cn.db.SetState(a.ID, taps.InProgress)
			}
		}
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.InProgress || d.State == taps.Done {
				cn.db.SetState(d.ID, taps.NotStarted)
				thsAffected[d.ID] = true
				stkItersAffected[d.Owner.Email+":"+d.Stks[d.Owner.Email].Iter] = true
			}
		}
		thsAffected[thread] = true
		stkItersAffected[th.Owner.Email+":"+th.Stks[th.Owner.Email].Iter] = true
	case taps.InProgress:
		for _, a := range cn.db.GetThreadAns(thread) {
			if a.State != taps.InProgress {
				cn.db.SetState(a.ID, taps.InProgress)
				stkItersAffected[a.Owner.Email+":"+a.Stks[a.Owner.Email].Iter] = true
			}
		}
		thsAffected[thread] = true
		stkItersAffected[th.Owner.Email+":"+th.Stks[th.Owner.Email].Iter] = true
	case taps.Done:
		for _, a := range cn.db.GetThreadAns(thread) {
			if a.State == taps.NotStarted {
				cn.db.SetState(a.ID, taps.InProgress)
				stkItersAffected[a.Owner.Email+":"+a.Stks[a.Owner.Email].Iter] = true
			}
		}
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.NotStarted || d.State == taps.InProgress {
				cn.db.SetState(d.ID, taps.Done)
				thsAffected[d.ID] = true
				stkItersAffected[d.Owner.Email+":"+d.Stks[d.Owner.Email].Iter] = true
			}
		}
		cn.db.SetPercentile(thread, 0)
		thsAffected[thread] = true
		stkItersAffected[th.Owner.Email+":"+th.Stks[th.Owner.Email].Iter] = true
	case taps.Closed:
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.NotStarted || d.State == taps.InProgress {
				cn.db.SetState(d.ID, taps.Closed)
				thsAffected[d.ID] = true
				stkItersAffected[d.Owner.Email+":"+d.Stks[d.Owner.Email].Iter] = true
			}

		}
		cn.db.SetPercentile(thread, 0)
		thsAffected[thread] = true
		stkItersAffected[th.Owner.Email+":"+th.Stks[th.Owner.Email].Iter] = true
	case taps.Archived:
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.NotStarted || d.State == taps.InProgress {
				cn.db.SetState(d.ID, taps.Archived)
				thsAffected[d.ID] = true
				stkItersAffected[d.Owner.Email+":"+d.Stks[d.Owner.Email].Iter] = true
			}
		}
		cn.db.SetPercentile(thread, 0)
		thsAffected[thread] = true
		stkItersAffected[th.Owner.Email+":"+th.Stks[th.Owner.Email].Iter] = true
	}
	for a := range thsAffected {
		cn.recalcAllCostTot(a)
		cn.recalcAllStkCosts(a)
	}
	for stkIter := range stkItersAffected {
		si := strings.Split(stkIter, ":")
		cn.recalcPri(si[0], si[1])
	}
}
