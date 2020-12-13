package tapfn

import "testing"

func TestRecalcPriOnThreadNew(t *testing.T) {
	cn, stks, _ := setupTest("1 stk")
	idA, err := cn.ThreadNew("A", stks["a"].Email, "2020-10 Oct", 1, nil, nil)
	if err != nil {
		t.Errorf("Could not insert thread: %v", err)
	}
	th, err := cn.Thread(idA)
	if err != nil {
		t.Errorf("Could not get thread: %v", err)
	}
	if x, g := 0.0, th.Percentile; x != g {
		t.Fatalf("Expected percentile %v; got %v", x, g)
	}
	idB, err := cn.ThreadNew("B", stks["a"].Email, "2020-10 Oct", 1, nil, nil)
	if err != nil {
		t.Errorf("Could not insert thread: %v", err)
	}
	th, err = cn.Thread(idB)
	if err != nil {
		t.Errorf("Could not get thread: %v", err)
	}
	if x, g := 0.5, th.Percentile; x != g {
		t.Fatalf("Expected percentile %v; got %v", x, g)
	}
}

func TestRecalcPriWithFullTree(t *testing.T) {
	cn, stks, _ := setupTest("1 stk")
	iter := "2020-10 Oct"
	stk := stks["a"].Email
	ths := map[string]int64{}
	type thTest struct {
		name   string
		parent string
		cost   int
		pct    float64
	}
	matrix := []thTest{
		thTest{
			name: "A",
			cost: 1,
			pct:  0.0,
		},
		thTest{
			name: "B",
			cost: 1,
			pct:  0.5,
		},
		thTest{
			name:   "A.A",
			parent: "A",
			cost:   1,
			pct:    0.0,
		},
		thTest{
			name:   "A.A.A",
			parent: "A.A",
			cost:   1,
			pct:    0.0,
		},
	}
	for _, m := range matrix {
		p := []int64{}
		if m.parent != "" {
			p = append(p, ths[m.parent])
		}
		id, err := cn.ThreadNew(m.name, stk, iter, m.cost, p, nil)
		if err != nil {
			t.Fatalf("Could not insert thread: %v", err)
		}
		ths[m.name] = id
		th, err := cn.Thread(id)
		if err != nil {
			t.Fatalf("Could not get thread: %v", err)
		}
		if x, g := m.pct, th.Percentile; x != g {
			t.Fatalf("Thread %v: expected percentile %v; got %v", m.name, x, g)
		}
	}
}

func TestRecalcPriOnSetOwner(t *testing.T) {
	cn, stks, ths := setupTest("2 stk, 3 th on 1")
	type mType struct {
		name string
		pct  float64
	}
	matrix := []mType{
		mType{name: "A", pct: 0.0},
		mType{name: "B", pct: 1.0 / 3},
		mType{name: "C", pct: 2.0 / 3},
	}
	for _, m := range matrix {
		th, err := cn.Thread(ths[m.name].ID)
		if err != nil {
			t.Fatalf("Could not get thread: %v", err)
		}
		if x, g := m.pct, th.Percentile; x != g {
			t.Fatalf("Thread %v before setting owner: Expected percentile %v; got %v", m.name, x, g)
		}
	}
	cn.ThreadSetOwner(ths["B"].ID, stks["b"].Email)
	matrix = []mType{
		mType{name: "A", pct: 0.0},
		mType{name: "B", pct: 0.0},
		mType{name: "C", pct: 1.0 / 2},
	}
	for _, m := range matrix {
		th, err := cn.Thread(ths[m.name].ID)
		if err != nil {
			t.Fatalf("Could not get thread: %v", err)
		}
		if x, g := m.pct, th.Percentile; x != g {
			t.Fatalf("Thread %v after setting owner: Expected percentile %v; got %v", m.name, x, g)
		}
	}
}
