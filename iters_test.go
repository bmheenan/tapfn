package tapfn

import "testing"

func TestGetItersForStk(t *testing.T) {
	cn, stks, _ := setupTest("s team")
	iters := cn.ItersByStk(stks["aa"].Email)
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
	cn, _, ths := setupTest("s team")
	iters := cn.ItersByParent(ths["A"].ID)
	expected := []string{"2020 Q4"}
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
