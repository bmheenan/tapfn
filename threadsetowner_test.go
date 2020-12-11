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
		t.Fatalf("Expected cost %v; got %v", x, g)
	}
}

func TestThreadSetOwnerAlreadyStk(t *testing.T) {
	cn, stks, ths := setupTest("2 stks w 1 th")
	cn.ThreadAddStk(ths["A"].ID, stks["b"].Email)
	cn.ThreadSetOwner(ths["A"].ID, stks["b"].Email)
	res := cn.ThreadrowsByStkIter(stks["b"].Email, "2020-12 Dec")
	if x, g := 1, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
	if x, g := "A", res[0].Name; x != g {
		t.Fatalf("Expected name %v; got %v", x, g)
	}
}
