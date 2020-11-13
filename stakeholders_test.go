package tapfn

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestNewAndGetStk(t *testing.T) {
	cn, errInit := Init(getTestCredentials())
	if errInit != nil {
		t.Errorf("Could not init: %v", errInit)
		return
	}
	errClear := cn.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Could not clear domain: %v", errClear)
		return
	}
	errN := cn.NewStk("a@example.com", "A", "A", "#ffffff", "#000000", taps.Monthly, []string{})
	if errN != nil {
		t.Errorf("Could not insert new stakeholder: %v", errN)
		return
	}
	stk, errStk := cn.GetStk("a@example.com")
	if errStk != nil {
		t.Errorf("Could not get stakeholder: %v", errStk)
		return
	}
	if stk.Name != "A" || stk.Email != "a@example.com" {
		t.Errorf(
			"Expected the stakeholder to have email:a@example.com, name:A. Got email:%v, name:%v",
			stk.Email,
			stk.Name,
		)
		return
	}
}

func TestNewWithParentAndGetStk(t *testing.T) {
	cn, errInit := Init(getTestCredentials())
	if errInit != nil {
		t.Errorf("Could not init: %v", errInit)
		return
	}
	errClear := cn.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Could not clear domain: %v", errClear)
		return
	}
	errN := cn.NewStk("a@example.com", "A", "A", "#ffffff", "#000000", taps.Monthly, []string{})
	if errN != nil {
		t.Errorf("Could not insert new stakeholder: %v", errN)
		return
	}
	errCh := cn.NewStk("b@example.com", "B", "B", "#ffffff", "#000000", taps.Biweekly, []string{"a@example.com"})
	if errCh != nil {
		t.Errorf("Could not insert new stakeholder child: %v", errCh)
		return
	}
	stkA, errA := cn.GetStk("a@example.com")
	if errA != nil {
		t.Errorf("Could not get stakeholder A: %v", errA)
		return
	}
	if stkA.Name != "A" || stkA.Email != "a@example.com" {
		t.Errorf(
			"Expected the stakeholder to have email:a@example.com, name:A. Got email:%v, name:%v",
			stkA.Email,
			stkA.Name,
		)
		return
	}
	stkB, errB := cn.GetStk("b@example.com")
	if errB != nil {
		t.Errorf("Could not get stakeholder B: %v", errB)
		return
	}
	if stkB.Name != "B" || stkB.Email != "b@example.com" {
		t.Errorf(
			"Expected the stakeholder to have email:b@example.com, name:B. Got email:%v, name:%v",
			stkB.Email,
			stkB.Name,
		)
		return
	}
}
