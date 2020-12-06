package tapfn

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bmheenan/taps"
)

func (cn *cnTapdb) ThreadNew(name, owner, iter string, cost int, parents, children []int64) (id int64, err error) {
	if name == "" || owner == "" || iter == "" || cost < 0 {
		return 0, errors.New("Name, owner, and iter must be non-blank; cost must be > 0")
	}
	oParts := strings.Split(owner, "@")
	if len(oParts) != 2 {
		return 0, fmt.Errorf("Owner email %v was invalid", owner)
	}
	id = cn.db.NewThread(name, oParts[1], owner, iter, string(taps.NotStarted), 1, cost)
	cn.ThreadAddStk(id, owner)
	// TODO set percentile
	for _, p := range parents {
		cn.ThreadLink(p, id)
	}
	for _, c := range children {
		cn.ThreadLink(id, c)
	}
	return id, nil
}
