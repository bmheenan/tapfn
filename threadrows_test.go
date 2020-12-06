package tapfn

import "testing"

func TestThreadrowsByStk(t *testing.T) {
	cn, stks, _ := setupTest("big tree")
	res := cn.ThreadrowsByStkIter(stks["a"].Email, "2020-10 Oct")
	if x, g := 2, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
}

func TestThreadrowsByParent(t *testing.T) {
	cn, _, ths := setupTest("big tree")
	res := cn.ThreadrowsByParentIter(ths["A"].ID, "2020-10 Oct")
	if x, g := 3, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
}

func TestThreadrowsByChild(t *testing.T) {
	cn, _, ths := setupTest("big tree")
	res := cn.ThreadrowsByChild(ths["AB"].ID)
	if x, g := 1, len(res); x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
	if x, g := "A", res[0].Name; x != g {
		t.Fatalf("Expected length %v; got %v", x, g)
	}
}
