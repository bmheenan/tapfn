package tapfn

import "testing"

func TestThreadAddAndRmStk(t *testing.T) {
	cn, stks, ths := setupTest("s team")
	thsAQ4 := cn.ThreadrowsByStkIter(stks["a"].Email, "2020 Q4")
	if x, g := 1, len(thsAQ4); x != g {
		t.Fatalf("A's Q4 list: Expected length %v; got %v", x, g)
	}
	if x, g := 15, thsAQ4[0].Cost; x != g {
		t.Fatalf("A's Q4 list: Expected cost %v; got %v", x, g)
	}
	cn.ThreadAddStk(ths["B"].ID, stks["a"].Email)
	thsAQ4 = cn.ThreadrowsByStkIter(stks["a"].Email, "2020 Q4")
	if x, g := 2, len(thsAQ4); x != g {
		t.Fatalf("A's Q4 list (after adding stk): Expected length %v; got %v", x, g)
	}
	if x, g := 10, thsAQ4[1].Cost; x != g {
		t.Fatalf("A's Q4 list (after adding stk): Expected cost %v; got %v", x, g)
	}
	cn.ThreadRemoveStk(ths["B"].ID, stks["a"].Email)
	thsAQ4 = cn.ThreadrowsByStkIter(stks["a"].Email, "2020 Q4")
	if x, g := 1, len(thsAQ4); x != g {
		t.Fatalf("A's Q4 list: Expected length %v; got %v", x, g)
	}
	if x, g := 15, thsAQ4[0].Cost; x != g {
		t.Fatalf("A's Q4 list: Expected cost %v; got %v", x, g)
	}
}
