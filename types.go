package tapfn

import (
	"errors"
	"time"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

// TapController provides fucntions for reading and changing application state
type TapController interface {

	// ClearDomain deletes all information within `domain` from the database
	DomainClear(domain string)

	// GetStk gets the information for the stakeholder with the given `email`
	Stk(email string) (taps.Stakeholder, error)

	// NewStk creates a new stakeholder with the given information. `email` must be unique. `name` is the display name.
	// `abbrev` should be a max 3 letter abbrevation of the name. `colorf` and `colorb` are the foreground and
	// background colors of the stakeholder's icon. `cadence` specifies how this stakholder tracks iterations. `parents`
	// nests this stakeholder under 0 or more existing stakeholders
	StkNew(email, name, abbrev, colorf, colorb string, cadence taps.Cadence, parents []string) error

	// GetStksForDomain returns a hierarchical view of all stakeholder in `domain`
	StksByDomain(domain string) []taps.StkInHier

	// GetItersForStk returns all iterations relevant to the given stakeholder `stk`, including those with a thread that
	// `stk` is a stakeholder for, plus Inbox, the current iteration, the next one, and Backlog
	ItersByStk(stk string) []string

	// GetItersForStk returns all iterations relevant to the given thread `parent`, including those with a thread
	// that's a child of `parent`, plus Inbox, the current iteration, the next one, and Backlog. They will be in the
	// cadence of `parent`'s owner
	ItersByParent(parent int64) []string

	// Thread returns the info for the given thread
	Thread(id int64) (th taps.Thread, err error)

	// ThreadNew creates a new thread with the given information. `name` is the name of the thread. `owner` is the email
	// of an existing stakeholder. `iter` is its iteration. `cost` is the direct cost of the thread. `parents` and
	// `children` nest this thread under existing ones, or next existing threads under this one
	ThreadNew(name, owner, iter string, cost int, parents, children []int64) (id int64, err error)

	// ThreadAddStk makes `stk` a stakeholder of `thread`, if not already
	ThreadAddStk(id int64, stk string)

	// ThreadRemoveStk makes `stk` no longer a stakeholder of `thread`, if it was
	ThreadRemoveStk(thread int64, stk string)

	// ThreadLink makes `parent` a parent of `child`. It returns an ErrNotFound if either does not exist.
	ThreadLink(parent, child int64) error

	// ThreadUnlink makes `parent` no longer a parent of `child`, if it was.
	ThreadUnlink(parent, child int64)

	// ThreadMoveForStk moves `thread` within its iteration as seen by stakeholder `stk`. Based on `moveTo`, it will be
	// moved to the start or end of the iteration, or right before `reference`
	ThreadMoveForStk(thread, reference int64, stkE string, moveTo MoveTo)

	// ThreadMoveForParent moves `thread` within its iteration as seen by parent `parent`. Based on `moveTo`, it will be
	// moved to the start or end of the iteration, or right before `reference`
	ThreadMoveForParent(thread, reference, parent int64, moveTo MoveTo)

	// ThreadrowsByStkIter returns all threadrows in hierarchical format for the given stakeholder `stk` and iteration
	// `iter`
	ThreadrowsByStkIter(stk, iter string) []taps.Threadrow

	// ThreadrowsByParentIter returns all threadrows in hierarchical format for the given parent `parent` and iteration
	// `iter`
	ThreadrowsByParentIter(parent int64, iter string) []taps.Threadrow

	// ThreadrowsByChild returns all threadrows (in a flat list) that are parents of the given `child`
	ThreadrowsByChild(child int64) []taps.Threadrow

	// SetThreadIter moves `thread` and all descendants in the same iteration to iteration `iter`
	ThreadSetIter(thread int64, iter string)

	// ThreadSetName sets the name of `thread` to `name`
	ThreadSetName(thread int64, name string)

	// ThreadSetDesc sets the description of `thread` to `desc`
	ThreadSetDesc(thread int64, desc string)

	// ThreadSetCost sets the direct cost of `thread` to `cost`, and recalculates total ancestor and stakeholder costs
	ThreadSetCost(thread int64, cost int)

	// ThreadSetState sets `thread` to `state`
	ThreadSetState(thread int64, state taps.State)

	// ThreadSetOwner sets `thread` to have owner `owner`
	ThreadSetOwner(thread int64, owner string)
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

// ErrWouldMakeLoop indicates that the items cannot be linked because loops are not allowed
var ErrWouldMakeLoop = errors.New("Cannot make a loop")

type cnTapdb struct {
	db           tapdb.DBInterface
	timeOverride time.Time // For testing
}

// MoveTo specifies the different anchors you can move a thread to within an iteration
type MoveTo int

const (
	// MoveToStart moves the thread to the beginning of the iteration, igoring the reference
	MoveToStart = iota
	// MoveToEnd moves the thread to the end of the iteration, ignoring the reference
	MoveToEnd
	// MoveBeforeRef moves the thread to right before the given reference
	MoveBeforeRef
)
