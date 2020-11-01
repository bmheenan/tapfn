package tapfn

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestNewAndGetPersonteam(t *testing.T) {
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
	cn.NewPersonteam("a@example.com", "A", "A", "#ffffff", "#000000", taps.Monthly, []string{})
	pt, errGet := cn.GetPersonteam("a@example.com")
	if errGet != nil {
		t.Errorf("Could not get personteam: %v", errGet)
		return
	}
	if pt.Name != "A" || pt.Email != "a@example.com" {
		t.Errorf(
			"Expected the personteam to have email:a@example.com, name:A. Got email:%v, name:%v",
			pt.Email,
			pt.Name,
		)
		return
	}
}
