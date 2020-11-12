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
	cn.NewStk("a@example.com", "A", "A", "#ffffff", "#000000", taps.Monthly, []string{})
	stk, errStk := cn.GetStk("a@example.com")
	if errStk != nil {
		t.Errorf("Could not get personteam: %v", errStk)
		return
	}
	if stk.Name != "A" || stk.Email != "a@example.com" {
		t.Errorf(
			"Expected the personteam to have email:a@example.com, name:A. Got email:%v, name:%v",
			pt.Email,
			pt.Name,
		)
		return
	}
}
