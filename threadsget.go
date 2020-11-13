package tapfn

import (
	"fmt"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) GetThread(id int64) (th *taps.Thread, err error) {
	th, err = cn.db.GetThread(id)
	if err != nil {
		err = fmt.Errorf("Could not get thread %v: %v", id, err)
	}
	return
}
