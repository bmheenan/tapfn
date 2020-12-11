package tapfn

import (
	"fmt"
	"testing"
)

func TestThreadSetOwner(t *testing.T) {
	cn, stks, ths := setupTest("big tree")
	_, err := cn.ThreadNew("AAB", stks["b"].Email, "2020-10 Oct", 3, []int64{ths["AA"].ID}, []int64{})
	if err != nil {
		t.Fatalf(fmt.Sprintf("Could not insert new thread: %v", err))
	}
	cn.ThreadSetOwner(ths["AA"].ID, stks["a"].Email)
	res := cn.ThreadrowsByStkIter(stks["a"].Email, "2020-10 Oct")
	if x, g := 2, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
	if x, g := 18, res[0].Cost; x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
}
