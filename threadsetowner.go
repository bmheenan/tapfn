package tapfn

import "fmt"

func (cn *cnTapdb) ThreadSetOwner(id int64, owner string) {
	th, err := cn.db.GetThread(id)
	if err != nil {
		panic(fmt.Sprintf("No thread with id %v", id))
	}
	po := th.Owner
	for _, d := range cn.db.GetThreadDes(id) {
		if d.Owner != po {
			continue
		}
		cn.ThreadAddStk(d.ID, owner)
		cn.db.SetOwner(d.ID, owner)
		cn.recalcAllStkCosts(d.ID)
	}
}
