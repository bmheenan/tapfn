package tapfn

import "testing"

func TestNewAndGetThread(t *testing.T) {
	cn, stks, iters, _, errSet := setupTest("1 stk")
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	id, errN := cn.NewThread("A", stks["a"].Email, iters[0], 1, nil, nil)
	if errN != nil {
		t.Errorf("Could not create thread A: %v", errN)
		return
	}
	th, errTh := cn.GetThread(id)
	if errTh != nil {
		t.Errorf("Could not get thread %v: %v", id, errTh)
		return
	}
	if th.Name != "A" {
		t.Errorf("Returned thread expected name A, got %v", th.Name)
		return
	}
}

func TestNewThreadWithParent(t *testing.T) {
	cn, stks, iters, ths, errSet := setupTest("1 th")
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	id, errN := cn.NewThread("AA", stks["a"].Email, iters[0], 1, []int64{ths["A"].ID}, nil)
	if errN != nil {
		t.Errorf("Could not make thread with parent: %v", errN)
		return
	}
	th, errTh := cn.GetThread(id)
	if errTh != nil {
		t.Errorf("Could not get inserted thread %v: %v", id, errTh)
		return
	}
	if th.Name != "AA" {
		t.Errorf("Expected thread name AA; got %v", th.Name)
		return
	}
	par, errP := cn.GetThread(ths["A"].ID)
	if errP != nil {
		t.Errorf("Could not get parent thread: %v", errP)
		return
	}
	if par.CostTot != 2 {
		t.Errorf("Expected parent total cost 2; got %v", par.CostTot)
		return
	}
}
