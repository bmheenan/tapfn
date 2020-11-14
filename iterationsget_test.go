package tapfn

import "testing"

func TestGetItersForStk(t *testing.T) {
	cn, stks, _, errSet := setupTest("s team")
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	iters, errI := cn.GetItersForStk(stks["aa"].Email)
	if errI != nil {
		t.Errorf("Could not get iterations for %v: %v", stks["aa"], errI)
		return
	}
	expected := []string{"Inbox", "2020-10 Oct", "2020-11 Nov", "2020-12 Dec", "Backlog"}
	if len(iters) != len(expected) {
		t.Errorf("Expected length %v; got %v", len(expected), len(iters))
		return
	}
	for i := range expected {
		if iters[i] != expected[i] {
			t.Errorf("Expected %v at indoex %v; got %v", expected[i], i, iters[i])
		}
	}
}

func TestGetItersForParent(t *testing.T) {
	cn, _, ths, errSet := setupTest("s team")
	if errSet != nil {
		t.Errorf("Could not set up test: %v", errSet)
		return
	}
	iters, errI := cn.GetItersForParent(ths["A"].ID)
	if errI != nil {
		t.Errorf("Could not get iterations for %v: %v", ths["A"].Name, errI)
		return
	}
	expected := []string{"Inbox", "2020 Q4", "2021 Q1", "Backlog"}
	if len(iters) != len(expected) {
		t.Errorf("Expected length %v; got %v", len(expected), len(iters))
		return
	}
	for i := range expected {
		if iters[i] != expected[i] {
			t.Errorf("Expected %v at indoex %v; got %v", expected[i], i, iters[i])
		}
	}
}
