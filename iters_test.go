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

func TestIterOptionsFuture(t *testing.T) {
	cn, _, ths := setupTest("2 stks w 1 th")
	iters := cn.IterOptions(ths["A"].ID)
	if x, g := 6, len(iters); x != g {
		t.Fatalf("Expected %v iterations; got %v", x, g)
	}
	if x, g := "2020-10 Oct", iters[1]; x != g {
		t.Fatalf("Expected 1st iter to be %v; got %v", x, g)
	}
	if x, g := "2021-01 Jan", iters[4]; x != g {
		t.Fatalf("Expected 1st iter to be %v; got %v", x, g)
	}
}

func TestIterOptionsPast(t *testing.T) {
	cn, _, ths := setupTest("th in past")
	iters := cn.IterOptions(ths["A"].ID)
	if x, g := 6, len(iters); x != g {
		t.Fatalf("Expected %v iterations; got %v", x, g)
	}
	if x, g := "2020-08 Aug", iters[1]; x != g {
		t.Fatalf("Expected 1st iter to be %v; got %v", x, g)
	}
	if x, g := "2020-11 Nov", iters[4]; x != g {
		t.Fatalf("Expected 1st iter to be %v; got %v", x, g)
	}
}
