package tapfn

import (
	"errors"
	"fmt"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) Thread(id int64) (th taps.Thread, err error) {
	thp, err := cn.db.GetThread(id)
	if errors.Is(err, tapdb.ErrNotFound) {
		th = taps.Thread{}
		err = fmt.Errorf("No thread with id %d: %w", id, ErrNotFound)
		return
	}
	if err != nil {
		panic(fmt.Sprintf("Unexpected error: %v", err))
	}
	th = *thp
	return
}
