package tapfn

import (
	"math"

	"github.com/bmheenan/taps"
)

type threadWork struct {
	cost int
	ord  int
}

func (cn *cnTapdb) recalcPri(stk string, iter string) {
	ths := cn.db.GetThreadsByStkIter(stk, iter)
	tws := map[int64]threadWork{}
	totCost := 0
	for _, th := range ths {
		if th.Owner.Email == stk && (th.State == taps.NotStarted || th.State == taps.InProgress) {
			tws[th.ID] = threadWork{
				cost: th.CostDir,
				ord:  th.Stks[stk].Ord,
			}
			totCost += th.CostDir
		}
	}
	for id, focus := range tws {
		des := cn.db.GetThreadDes(id)
		var (
			totEarlier int
			percentile float64
		)
		for i, tw := range tws {
			if tw.ord > focus.ord {
				continue
			}
			if _, ok := des[i]; ok {
				continue
			}
			totEarlier += tw.cost
		}
		percentile = float64(totEarlier) / float64(totCost)
		if math.IsNaN(percentile) {
			percentile = 0
		}
		cn.db.SetPercentile(id, percentile)
	}
}
