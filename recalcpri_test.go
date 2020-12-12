package tapfn

import "testing"

func TestRecalcPriOnThreadNew(t *testing.T) {
	cn, stks, _ := setupTest("1 stk")
	idA, err := cn.ThreadNew("A", stks["a"].Email, "2020-10 Oct", 1, nil, nil)
	if err != nil {
		t.Errorf("Could not insert thread: %v", err)
	}
	th, err := cn.Thread(idA)
	if err != nil {
		t.Errorf("Could not get thread: %v", err)
	}
	if x, g := 0.0, th.Percentile; x != g {
		t.Fatalf("Expected percentile %v; got %v", x, g)
	}
	idB, err := cn.ThreadNew("B", stks["a"].Email, "2020-10 Oct", 1, nil, nil)
	if err != nil {
		t.Errorf("Could not insert thread: %v", err)
	}
	th, err = cn.Thread(idB)
	if err != nil {
		t.Errorf("Could not get thread: %v", err)
	}
	if x, g := 0.5, th.Percentile; x != g {
		t.Fatalf("Expected percentile %v; got %v", x, g)
	}
}
