package tapfn

import "testing"

func TestNewChildInDifferentIter(t *testing.T) {
	cn, stks, ths, err := setupTest("1 stk w ths")
	if err != nil {
		t.Fatalf("Could not set up test: %v", err)
	}
	id, err := cn.NewThread("New thread", stks["a"].Email, "2020-10 Oct", 1, []int64{ths["A"].ID}, []int64{})
	if err != nil {
		t.Fatalf("Could not create thread: %v", err)
	}
	res, err := cn.GetThreadrowsByParentIter(ths["A"].ID, "2020-10 Oct")
	if err != nil {
		t.Fatalf("Could not get threadrows: %v", err)
	}
	th, err := cn.GetThread(id)
	expect("result iteration", "2020-10 Oct", th.Iter, t)
	expect("result iteration for parent", "2020-10 Oct", th.Parents[ths["A"].ID].Iter, t)
	expect("results length", 1, len(res), t)
}
