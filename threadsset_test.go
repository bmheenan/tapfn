package tapfn

import "testing"

func TestSetIter(t *testing.T) {
	cn, _, ths, err := setupTest("s team")
	if err != nil {
		t.Errorf("Could not set up test: %v", err)
		return
	}
	err = cn.SetThreadIter(ths["A"].ID, "2021 Q1")
	if err != nil {
		t.Errorf("Could not ser thread iteration: %v", err)
		return
	}
	th, err := cn.GetThread(ths["A"].ID)
	if err != nil {
		t.Errorf("Could not get thread: %v", err)
		return
	}
	if x := "2021 Q1"; th.Iter != x {
		t.Errorf("Expected iter %v; got %v", x, th.Iter)
		return
	}
}

func TestSetIterChildren(t *testing.T) {
	cn, _, ths, err := setupTest("1 stk w ths")
	if err != nil {
		t.Errorf("Could not set up test: %v", err)
		return
	}
	err = cn.SetThreadIter(ths["A"].ID, "2021-01 Jan")
	if err != nil {
		t.Errorf("Could not set thread iteration: %v", err)
		return
	}
	for _, n := range []string{"A", "AA", "AB", "AC"} {
		th, err := cn.GetThread(ths[n].ID)
		if err != nil {
			t.Errorf("Could not get thread: %v", err)
			return
		}
		if x := "2021-01 Jan"; th.Iter != x {
			t.Errorf("Expected iter %v; got %v", x, th.Iter)
			return
		}
	}

}

func TestSetIterStkPa(t *testing.T) {
	cn, stks, ths, err := setupTest("1 stk w ths")
	if err != nil {
		t.Errorf("Could not set up test: %v", err)
		return
	}
	err = cn.SetThreadIter(ths["A"].ID, "2021-01 Jan")
	if err != nil {
		t.Errorf("Could not set thread iteration: %v", err)
		return
	}
	res, err := cn.GetThreadrowsByStkIter(stks["a"].Email, "2021-01 Jan")
	if err != nil {
		t.Fatalf("Could not get threadrows: %v", err)
	}
	if x, g := 1, len(res); x != g {
		t.Fatalf("Expected length %d; got %d", x, g)
	}
	if x, g := 3, len(res[0].Children); x != g {
		t.Fatalf("Expected length %d; got %d", x, g)
	}
	res, err = cn.GetThreadrowsByParentIter(ths["A"].ID, "2021-01 Jan")
	if err != nil {
		t.Fatalf("Could not get threadrows for parent")
	}
	if x, g := 3, len(res); x != g {
		t.Fatalf("Expected length %d; got %d", x, g)
	}
}
