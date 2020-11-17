package tapfn

import "testing"

func TestNewAndGetThread(t *testing.T) {
	cn, stks, _, errSet := setupTest("1 stk")
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	id, errN := cn.NewThread("A", stks["a"].Email, "2020 Oct", 1, nil, nil)
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
	cn, stks, ths, errSet := setupTest("1 th")
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	id, errN := cn.NewThread("AA", stks["a"].Email, "2020 Q1", 1, []int64{ths["A"].ID}, nil)
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

func TestMoveThreadForStk(t *testing.T) {
	cn, stks, ths, errSet := setupTest("s team")
	if errSet != nil {
		t.Errorf("Could not setup test: %v", errSet)
	}
	errM := cn.MoveThreadForStk(ths["AC"].ID, ths["AB"].ID, stks["ab"].Email)
	if errM != nil {
		t.Errorf("Could not move AC before AB: %v", errM)
		return
	}
	thrs, errThs := cn.GetThreadrowsByStkIter(stks["ab"].Email, "2020-10 Oct")
	if errThs != nil {
		t.Errorf("Could not get threadrows for %v 2020 Oct: %v", stks["ab"].Email, errThs)
		return
	}
	if len(thrs) != 2 {
		t.Errorf("Expected 2 results; got %v", len(thrs))
		return
	}
	if thrs[0].Name != "AC" {
		t.Errorf("Expected first threadrow to be AC; was %v", thrs[0].Name)
		return
	}
	if thrs[1].Name != "AB" {
		t.Errorf("Expected second threadrow to be AB; was %v", thrs[0].Name)
		return
	}
}

func TestMoveThreadToEndForStk(t *testing.T) {
	cn, stks, ths, errSet := setupTest("s team")
	if errSet != nil {
		t.Errorf("Could not setup test: %v", errSet)
	}
	errM := cn.MoveThreadForStk(ths["AB"].ID, 0, stks["ab"].Email)
	if errM != nil {
		t.Errorf("Could not move AC before AB: %v", errM)
		return
	}
	thrs, errThs := cn.GetThreadrowsByStkIter(stks["ab"].Email, "2020-10 Oct")
	if errThs != nil {
		t.Errorf("Could not get threadrows for %v 2020 Oct: %v", stks["ab"].Email, errThs)
		return
	}
	if len(thrs) != 2 {
		t.Errorf("Expected 2 results; got %v", len(thrs))
		return
	}
	if thrs[0].Name != "AC" {
		t.Errorf("Expected first threadrow to be AC; was %v", thrs[0].Name)
		return
	}
	if thrs[1].Name != "AB" {
		t.Errorf("Expected second threadrow to be AB; was %v", thrs[0].Name)
		return
	}
}

func TestMoveThreadForParent(t *testing.T) {
	cn, stks, ths, errSet := setupTest("s team")
	if errSet != nil {
		t.Errorf("Could not setup test: %v", errSet)
	}
	errM := cn.MoveThreadForParent(ths["AC"].ID, ths["AA"].ID, ths["A"].ID)
	if errM != nil {
		t.Errorf("Could not move AC before AB: %v", errM)
		return
	}
	thrs, errThs := cn.GetThreadrowsByParentIter(ths["A"].ID, "2020 Q4")
	if errThs != nil {
		t.Errorf("Could not get threadrows for %v Q4: %v", stks["ab"].Email, errThs)
		return
	}
	if len(thrs) != 3 {
		t.Errorf("Expected 2 results; got %v", len(thrs))
		return
	}
	if thrs[0].Name != "AC" {
		t.Errorf("Expected first threadrow to be AC; was %v", thrs[0].Name)
		return
	}
	if thrs[1].Name != "AA" {
		t.Errorf("Expected second threadrow to be AA; was %v", thrs[1].Name)
		return
	}
	if thrs[2].Name != "AB" {
		t.Errorf("Expected third threadrow to be AB; was %v", thrs[2].Name)
		return
	}
}

func TestMoveThreadToEndForParent(t *testing.T) {
	cn, stks, ths, errSet := setupTest("s team")
	if errSet != nil {
		t.Errorf("Could not setup test: %v", errSet)
		return
	}
	errM := cn.MoveThreadForParent(ths["AB"].ID, 0, ths["A"].ID)
	if errM != nil {
		t.Errorf("Could not move AB to end of the iteration: %v", errM)
		return
	}
	thrs, errThs := cn.GetThreadrowsByParentIter(ths["A"].ID, "2020 Q4")
	if errThs != nil {
		t.Errorf("Could not get threadrows for %v Q4: %v", stks["ab"].Email, errThs)
		return
	}
	if len(thrs) != 3 {
		t.Errorf("Expected 2 results; got %v", len(thrs))
		return
	}
	if thrs[0].Name != "AA" {
		t.Errorf("Expected first threadrow to be AA; was %v", thrs[0].Name)
		return
	}
	if thrs[1].Name != "AC" {
		t.Errorf("Expected second threadrow to be AC; was %v", thrs[1].Name)
		return
	}
	if thrs[2].Name != "AB" {
		t.Errorf("Expected third threadrow to be AB; was %v", thrs[2].Name)
		return
	}
}
