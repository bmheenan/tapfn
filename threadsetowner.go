package tapfn

import (
	"fmt"
	"strings"
)

func (cn *cnTapdb) ThreadSetOwner(id int64, owner string) {
	th, err := cn.db.GetThread(id)
	if err != nil {
		panic(fmt.Sprintf("No thread with id %v", id))
	}
	po := th.Owner
	stkItersAffected := map[string]bool{}
	for _, d := range cn.db.GetThreadDes(id) {
		if d.Owner != po {
			continue
		}
		cn.ThreadAddStk(d.ID, owner)
		cn.db.SetOwner(d.ID, owner)
		cn.recalcAllStkCosts(d.ID)
		stkItersAffected[d.Owner.Email+":"+d.Stks[d.Owner.Email].Iter] = true
		stkItersAffected[owner+":"+d.Stks[owner].Iter] = true
	}
	for stkIter := range stkItersAffected {
		si := strings.Split(stkIter, ":")
		cn.recalcPri(si[0], si[1])
	}
}
