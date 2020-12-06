package tapfn

import "testing"

func TestThreadSetIter(t *testing.T) {
	cn, stks, ths := setupTest("big tree")
	_, err := cn.ThreadNew("D", stks["a"].Email, "2020-11 Nov", 5, []int64{}, []int64{})
	if err != nil {
		t.Fatalf("Could not insert new thread: %v", err)
	}
	cn.ThreadSetIter(ths["A"].ID, "2020-11 Nov")
	res := cn.ThreadrowsByStkIter(stks["a"].Email, "2020-11 Nov")
	if x, g := 2, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
}
