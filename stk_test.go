package tapfn

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestNewAndGetStk(t *testing.T) {
	cn, _, _ := setupTest("blank")
	err := cn.StkNew("a@example.com", "A", "A", "#ffffff", "#000000", taps.Monthly, []string{})
	if err != nil {
		t.Fatalf("Could not insert new stakeholder: %v", err)
	}
	stk, errStk := cn.Stk("a@example.com")
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
	cn, _, _ := setupTest("blank")
	cn.DomainClear("example.com")
	err := cn.StkNew("a@example.com", "A", "A", "#ffffff", "#000000", taps.Monthly, []string{})
	if err != nil {
		t.Fatalf("Could not insert new stakeholder: %v", err)
	}
	err = cn.StkNew("b@example.com", "B", "B", "#ffffff", "#000000", taps.Biweekly, []string{"a@example.com"})
	if err != nil {
		t.Fatalf("Could not insert new stakeholder child: %v", err)
	}
	stkA, err := cn.Stk("a@example.com")
	if err != nil {
		t.Fatalf("Could not get stakeholder A: %v", err)
	}
	if stkA.Name != "A" || stkA.Email != "a@example.com" {
		t.Errorf(
			"Expected the stakeholder to have email:a@example.com, name:A. Got email:%v, name:%v",
			stkA.Email,
			stkA.Name,
		)
		return
	}
	stkB, err := cn.Stk("b@example.com")
	if err != nil {
		t.Fatalf("Could not get stakeholder B: %v", err)
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

func TestGetStksForDomain(t *testing.T) {
	cn, _, _ := setupTest("s team no ths")
	stks := cn.StksByDomain("example.com")
	if len(stks) != 2 {
		t.Errorf("Expected len 2; got %v", len(stks))
		return
	}
	if stks[0].Name != "Team a" {
		t.Errorf("Expected Team a; got %v", stks[0].Name)
		return
	}
	if stks[0].Members[0].Name != "Person aa" {
		t.Errorf("Expected Person aa; got %v", stks[0].Members[0].Name)
		return
	}
	if stks[0].Members[1].Name != "Person ab" {
		t.Errorf("Expected Person ab; got %v", stks[0].Members[1].Name)
		return
	}
	if stks[1].Name != "Team b" {
		t.Errorf("Expected Team b; got %v", stks[1].Name)
		return
	}
}
