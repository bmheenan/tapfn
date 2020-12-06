package tapfn

import "testing"

func TestNewAndGetThread(t *testing.T) {
	cn, stks, _ := setupTest("1 stk")
	id, err := cn.ThreadNew("A", stks["a"].Email, "2020-10 Oct", 1, nil, nil)
	if err != nil {
		t.Fatalf("Could not create thread A: %v", err)
	}
	th, err := cn.Thread(id)
	if err != nil {
		t.Fatalf("Could not get thread %v: %v", id, err)
	}
	if th.Name != "A" {
		t.Fatalf("Returned thread expected name A, got %v", th.Name)
	}
}

func TestNewThreadWithParent(t *testing.T) {
	cn, stks, ths := setupTest("1 th")
	id, errN := cn.ThreadNew("AA", stks["a"].Email, "2020 Q1", 1, []int64{ths["A"].ID}, nil)
	if errN != nil {
		t.Errorf("Could not make thread with parent: %v", errN)
		return
	}
	th, errTh := cn.Thread(id)
	if errTh != nil {
		t.Errorf("Could not get inserted thread %v: %v", id, errTh)
		return
	}
	if th.Name != "AA" {
		t.Errorf("Expected thread name AA; got %v", th.Name)
		return
	}
	par, errP := cn.Thread(ths["A"].ID)
	if errP != nil {
		t.Errorf("Could not get parent thread: %v", errP)
		return
	}
	if par.CostTot != 2 {
		t.Errorf("Expected parent total cost 2; got %v", par.CostTot)
		return
	}
}

func TestThreadUnlink(t *testing.T) {
	cn, _, ths := setupTest("s team")
	cn.ThreadUnlink(ths["A"].ID, ths["AC"].ID)
	th, err := cn.Thread(ths["AC"].ID)
	if err != nil {
		t.Errorf("Could not get thread AC: %v", err)
		return
	}
	if _, ok := th.Parents[ths["A"].ID]; ok {
		t.Errorf("Expected AC not to have parent A, but found it")
		return
	}
}

func TestStkCostInDiffIters(t *testing.T) {
	cn, stks, ths := setupTest("big tree")
	cn.ThreadSetIter(ths["AB"].ID, "2020-09 Sep")
	res := cn.ThreadrowsByStkIter(stks["a"].Email, "2020-10 Oct")
	if x, g := 5, res[0].Cost; x != g {
		t.Fatalf("Expected A to have cost %v; got %v", x, g)
	}
}
