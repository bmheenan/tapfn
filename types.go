package tapfn

import (
	"errors"
	"time"

	"github.com/bmheenan/tapdb"
	"github.com/bmheenan/taps"
)

// TapController provides fucntions for reading and changing application state
type TapController interface {

	// domain

	// ClearDomain deletes all information within `domain` from the database
	DomainClear(domain string)

	// stk

	// GetStk gets the information for the stakeholder with the given `email`
	Stk(email string) (taps.Stakeholder, error)

	// NewStk creates a new stakeholder with the given information. `email` must be unique. `name` is the display name.
	// `abbrev` should be a max 3 letter abbrevation of the name. `colorf` and `colorb` are the foreground and
	// background colors of the stakeholder's icon. `cadence` specifies how this stakholder tracks iterations. `parents`
	// nests this stakeholder under 0 or more existing stakeholders
	StkNew(email, name, abbrev, colorf, colorb string, cadence taps.Cadence, parents []string) error

	// GetStksForDomain returns a hierarchical view of all stakeholder in `domain`
	StksByDomain(domain string) []taps.StkInHier

	// iters

	// GetItersForStk returns all iterations relevant to the given stakeholder `stk`, including those with a thread that
	// `stk` is a stakeholder for, plus Inbox, the current iteration, the next one, and Backlog
	ItersForStk(stk string) []string

	// GetItersForStk returns all iterations relevant to the given thread `parent`, including those with a thread
	// that's a child of `parent`, plus Inbox, the current iteration, the next one, and Backlog. They will be in the
	// cadence of `parent`'s owner
	ItersForParent(parent int64) []string

	// thread

	// Thread returns the info for the given thread
	Thread(id int64) (th taps.Thread, err error)

	// ThreadNew creates a new thread with the given information. `name` is the name of the thread. `owner` is the email
	// of an existing stakeholder. `iter` is its iteration. `cost` is the direct cost of the thread. `parents` and
	// `children` nest this thread under existing ones, or next existing threads under this one
	ThreadNew(name, owner, iter string, cost int, parents, children []int64) (id int64, err error)

	// ThreadAddStk makes `stk` a stakholder of `thread`, if not already
	ThreadAddStk(id int64, stk string)

	/*
		// threadsdelete.go

		// DeleteThreadHierLinks removes all links from `child` to each of its parents that link it up to ancestor `anc`
		DeleteThreadHierLinks(anc, child int64) error

		// threadsget.go

		// GetThread returns the info for the given thread
		GetThread(id int64) (th taps.Thread, err error)

		// GetThreadrowsByStkIter gets all threads where `stk` is a stakeholder in iteration `iter`. They're returned as
		// threadrows scoped to `stk`, so they use `stk`'s ordering and only show costs for the pieces that `stk` owns
		// (including team members)
		GetThreadrowsByStkIter(stk, iter string) (ths []taps.Threadrow, err error)

		// GetThreadrowsByParentIter gets all threads that are children of `parent`, and recursively gets their children
		// until all descendants are fetched. They're returned as threadrows scoped to `parent`, so they use `parent`'s
		// order, and display their total cost
		GetThreadrowsByParentIter(parent int64, iter string) (ths []taps.Threadrow, err error)

		// GetThreadrowsByChild gets all threads that are direct parents of the given `child` thread
		GetThreadrowsByChild(child int64) (ths []taps.Threadrow, err error)

		// threadsnew.go

		// NewThread creates a new thread with the given information. `name` is the name of the thread. `owner` is the email
		// of an existing stakeholder. `iter` is its iteration. `cost` is the direct cost of the thread. `parents` and
		// `children` nest this thread under existing ones, or next existing threads under this one
		NewThread(name, owner, iter string, cost int, parents, children []int64) (int64, error)

		// NewThreadHierLink links a `parent` thread with a `child` in the hierarchy. You cannot create a loop
		NewThreadHierLink(parent, child int64) error

		// threadsmove.go

		// MoveThreadParent moves the thread with id `thread` in the context of the thread with id `parent`. `parent` must be a
		// parent of the other two threads, and the two threads must be in the same iteration
		// You can move the thread to the start or end of the iteration, or immediately before the reference thread, depending
		// on the value of `anchor`
		MoveThreadForParent(thread, reference, parent int64, moveTo MoveTo) error

		// MoveThreadStakeholder moves the thread with id `thread` within its iteration in the context of `stkE`.
		// `stakeholder` is a stakeholder of the thread. If moveTo = MoveToStart or MoveToEnd, the thread will be moved to the
		// start or end of the iteration. If moveTo = MoveBeforeRef, `reference` must be a thread with the same iteration and
		// stakeholder, and `thread` will be moved immediately before it.
		MoveThreadForStk(thread, reference int64, stk string, moveTo MoveTo) error

		// threadsset.go

		// SetThreadIter moves `thread` and all descendants in the same iteration to iteration `iter`
		SetThreadIter(thread int64, iter string) (err error)
	*/
}

// ErrNotFound indicates that no matching record was found when querying
var ErrNotFound = errors.New("Not found")

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
