package tapfn

import "github.com/bmheenan/taps"

func (cn *cnTapdb) ThreadSetState(thread int64, state taps.State) {
	cn.db.SetState(thread, state)
	affected := map[int64]bool{}
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
			}
			affected[d.ID] = true
		}
	case taps.InProgress:
		for _, a := range cn.db.GetThreadAns(thread) {
			if a.State != taps.InProgress {
				cn.db.SetState(a.ID, taps.InProgress)
			}
		}
		affected[thread] = true
	case taps.Done:
		for _, a := range cn.db.GetThreadAns(thread) {
			if a.State == taps.NotStarted {
				cn.db.SetState(a.ID, taps.InProgress)
			}
		}
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.NotStarted || d.State == taps.InProgress {
				cn.db.SetState(d.ID, taps.Done)
			}
			affected[d.ID] = true
		}
	case taps.Closed:
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.NotStarted || d.State == taps.InProgress {
				cn.db.SetState(d.ID, taps.Closed)
			}
			affected[d.ID] = true
		}
	case taps.Archived:
		for _, d := range cn.db.GetThreadDes(thread) {
			if d.State == taps.NotStarted || d.State == taps.InProgress {
				cn.db.SetState(d.ID, taps.Archived)
			}
			affected[d.ID] = true
		}
	}
	for a := range affected {
		cn.recalcAllCostTot(a)
		cn.recalcAllStkCosts(a)
	}
}
