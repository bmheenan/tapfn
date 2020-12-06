package tapfn

import "testing"

func TestMoveThreadForStk(t *testing.T) {
	cn, stks, ths := setupTest("s team")
	cn.ThreadMoveForStk(ths["AC"].ID, ths["AB"].ID, stks["ab"].Email, MoveBeforeRef)
	thrs := cn.ThreadrowsByStkIter(stks["ab"].Email, "2020-10 Oct")
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
	cn, stks, ths := setupTest("s team")
	cn.ThreadMoveForStk(ths["AB"].ID, 0, stks["ab"].Email, MoveToEnd)
	thrs := cn.ThreadrowsByStkIter(stks["ab"].Email, "2020-10 Oct")
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

func TestMoveThreadToStartForStk(t *testing.T) {
	cn, stks, ths := setupTest("s team")
	cn.ThreadMoveForStk(ths["AC"].ID, 0, stks["ab"].Email, MoveToStart)
	thrs := cn.ThreadrowsByStkIter(stks["ab"].Email, "2020-10 Oct")
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
	cn, _, ths := setupTest("s team")
	cn.ThreadMoveForParent(ths["AC"].ID, ths["AA"].ID, ths["A"].ID, MoveBeforeRef)
	thrs := cn.ThreadrowsByParentIter(ths["A"].ID, "2020 Q4")
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
	cn, _, ths := setupTest("s team")
	cn.ThreadMoveForParent(ths["AB"].ID, 0, ths["A"].ID, MoveToEnd)
	thrs := cn.ThreadrowsByParentIter(ths["A"].ID, "2020 Q4")
	if len(thrs) != 3 {
		t.Errorf("Expected 3 results; got %v", len(thrs))
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

func TestMoveThreadToStartForParent(t *testing.T) {
	cn, _, ths := setupTest("s team")
	cn.ThreadMoveForParent(ths["AB"].ID, 0, ths["A"].ID, MoveToStart)
	thrs := cn.ThreadrowsByParentIter(ths["A"].ID, "2020 Q4")
	if len(thrs) != 3 {
		t.Errorf("Expected 3 results; got %v", len(thrs))
		return
	}
	if thrs[0].Name != "AB" {
		t.Errorf("Expected first threadrow to be AB; was %v", thrs[0].Name)
		return
	}
	if thrs[1].Name != "AA" {
		t.Errorf("Expected second threadrow to be AA; was %v", thrs[1].Name)
		return
	}
	if thrs[2].Name != "AC" {
		t.Errorf("Expected third threadrow to be AC; was %v", thrs[2].Name)
		return
	}
}
