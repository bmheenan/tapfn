package tapfn

import "fmt"

func (cn *cnTapdb) ThreadSetIter(thread int64, iter string) {
	th, err := cn.db.GetThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get thread: %v", err))
	}
	for _, dec := range cn.db.GetThreadDes(thread) {
		if iterResulting(dec.Iter, th.Owner.Cadence) != th.Iter {
			break
		}
		iter := iterResulting(iter, dec.Owner.Cadence)
		cn.db.SetIter(dec.ID, iter)
		for parent := range dec.Parents {
			pa, err := cn.db.GetThread(parent)
			if err != nil {
				panic(fmt.Sprintf("Could not get parent thread: %v", err))
			}
			iter := iterResulting(iter, pa.Owner.Cadence)
			var place MoveTo
			switch {
			case iter == dec.Parents[parent].Iter:
				break
			case iter < dec.Parents[parent].Iter:
				place = MoveToEnd
			case iter > dec.Parents[parent].Iter:
				place = MoveToStart
			}
			cn.db.SetIterForParent(parent, dec.ID, iter)
			cn.ThreadMoveForParent(dec.ID, 0, parent, place)
		}
		for stk := range dec.Stks {
			s, err := cn.Stk(stk)
			if err != nil {
				panic(fmt.Sprintf("Could not get stakeholder: %v", err))
			}
			iter := iterResulting(iter, s.Cadence)
			var place MoveTo
			switch {
			case iter == dec.Stks[stk].Iter:
				break
			case iter < dec.Stks[stk].Iter:
				place = MoveToEnd
			case iter > dec.Stks[stk].Iter:
				place = MoveToStart
			}
			cn.db.SetIterForStk(dec.ID, stk, iter)
			cn.ThreadMoveForStk(dec.ID, 0, stk, place)
		}
	}
}
