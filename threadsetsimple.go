package tapfn

import (
	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) ThreadSetName(thread int64, name string) {
	cn.db.SetName(thread, name)
}

func (cn *cnTapdb) ThreadSetDesc(thread int64, desc string) {
	cn.db.SetDesc(thread, desc)
}

func (cn *cnTapdb) ThreadSetState(thread int64, state taps.State) {
	cn.db.SetState(thread, state)
}
