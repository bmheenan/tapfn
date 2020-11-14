package tapfn

import (
	"errors"
	"time"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

// TapController provides fucntions for reading and changing application state
type TapController interface {

	// cleardomain.go

	// ClearDomain deletes all information within `domain` from the database
	ClearDomain(domain string) error

	// stakeholders.go

	// NewStk creates a new stakeholder with the given information. `email` must be unique. `name` is the display name.
	// `abbrev` should be a max 3 letter abbrevation of the name. `colorf` and `colorb` are the foreground and
	// background colors of the stakeholder's icon. `cadence` specifies how this stakholder tracks iterations. `parents`
	// nests this stakeholder under 0 or more existing stakeholders
	NewStk(email, name, abbrev, colorf, colorb string, cadence taps.Cadence, parents []string) error

	// GetStk gets the information for the stakeholder with the given `email`
	GetStk(email string) (*taps.Stakeholder, error)

	// iterationsget.go

	// GetItersForStk returns all iterations relevant to the given stakeholder `stk`, including those with a thread that
	// `stk` is a stakeholder for, plus Inbox, the current iteration, the next one, and Backlog
	GetItersForStk(stk string) (iters []string, err error)

	// GetItersForStk returns all iterations relevant to the given thread `parent`, including those with a thread
	// that's a child of `parent`, plus Inbox, the current iteration, the next one, and Backlog. They will be in the
	// cadence of `parent`'s owner
	GetItersForParent(parent int64) (iters []string, err error)

	// threadsget.go

	// GetThread returns the info for the given thread
	GetThread(id int64) (th *taps.Thread, err error)

	GetThreadrowsByStkIter(stk, iter string) (ths []*taps.Threadrow, err error)

	GetThreadrowsByParentIter(parent int64, iter string) (ths []*taps.Threadrow, err error)

	// threadsnew.go

	// NewThread creates a new thread with the given information. `name` is the name of the thread. `owner` is the email
	// of an existing stakeholder. `iter` is its iteration. `cost` is the direct cost of the thread. `parents` and
	// `children` nest this thread under existing ones, or next existing threads under this one
	NewThread(name, owner, iter string, cost int, parents, children []int64) (int64, error)

	// NewThreadHierLink links a `parent` thread with a `child` in the hierarchy. You cannot create a loop
	NewThreadHierLink(parent, child int64) error

	// threadsmove.go

	// MoveThreadByParent changes the order of a `thread` with respect to a `parent`. It doesn't affect the order for
	// other parents or any stakeholders. `thread` will be moved immediately before `reference`, or to the end of the
	// iteraton if `reference` == 0
	MoveThreadForParent(thread, reference, parent int64) error

	// MoveThreadByStk changes the order of a `thread` with respect to a stakeholder `stk`. It doesn't affect the order
	// for other stakeholders or any parents. `thread` will be moved immediately before `reference`, or to the end of
	// the iteraton if `reference` == 0
	MoveThreadForStk(thread, reference int64, stk string) error
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

type cnTapdb struct {
	db           tapdb.DBInterface
	timeOverride time.Time // For testing
}
