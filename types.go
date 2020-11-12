package tapfn

import (
	"errors"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

// TapController provides fucntions for reading and changing application state
type TapController interface {
	ClearDomain(domain string) error

	NewStk(email, name, abbrev, colorf, colorb string, itertiming taps.Cadence, parents []string) error
	GetStk(email string) (*taps.Stakeholder, error)

	NewThread(name, owner, iter string, cost int, parents, children []int64) (int64, error)
	LinkThreads(parent, child int64) error
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

type cnTapdb struct {
	db tapdb.DBInterface
}
