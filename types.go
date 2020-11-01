package tapfn

import (
	"errors"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

// TapController provides fucntions for reading and changing application state
type TapController interface {
	ClearDomain(domain string) error
	NewPersonteam(email, name, abbrev, colorf, colorb string, itertiming taps.IterTiming, parents []string) error
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

// ErrBadArgs indicates that the arguments given to the function were bad
var ErrBadArgs = errors.New("Bad arguments")

type cnTapdb struct {
	db tapdb.DBInterface
}
