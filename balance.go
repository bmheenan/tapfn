package tapfn

import (
	"math"
)

func (cn *cnTapdb) balanceParent(parent int64, iter string) {
	ths := cn.db.GetThreadsByParentIter(parent, iter)
	step := math.MaxInt32 / (len(ths) + 1)
	for i, th := range ths {
		cn.db.SetOrdForParent(th.ID, parent, step*(i+1))
	}
}

func (cn *cnTapdb) balanceStk(stk, iter string) {
	ths := cn.db.GetThreadsByStkIter(stk, iter)
	step := math.MaxInt32 / (len(ths) + 1)
	for i, th := range ths {
		cn.db.SetOrdForStk(th.ID, stk, step*(i+1))
	}
}
