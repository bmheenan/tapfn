package tapfn

import "testing"

func TestInit(t *testing.T) {
	_, err := Init(getTestCredentials())
	if err != nil {
		t.Errorf("Could not initialize TapController: %v", err)
		return
	}
}
