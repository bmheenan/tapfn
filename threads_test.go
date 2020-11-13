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
