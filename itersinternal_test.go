package tapfn

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestIterationQToM(t *testing.T) {
	i := iterResulting("2020 Q3", taps.Monthly)
	if i != "2020-09 Sep" {
		t.Fatalf("Expected 2020-09 Sep; got %v", i)
	}
}

func TestIterationMToQ(t *testing.T) {
	i := iterResulting("2020-02 Feb", taps.Quarterly)
	if i != "2020 Q1" {
		t.Fatalf("Expected 2020 Q1; got %v", i)
	}
}

func TestIterationMToY(t *testing.T) {
	i := iterResulting("2021-06 Jun", taps.Yearly)
	if i != "2021" {
		t.Fatalf("Expected 2021; got %v", i)
	}
}

func TestIterationYToM(t *testing.T) {
	i := iterResulting("2019", taps.Monthly)
	if i != "2019-12 Dec" {
		t.Fatalf("Expected 2019-12 Dec; got %v", i)
	}
}

func TestGetCadenceY(t *testing.T) {
	c, v := iterCadence("2020")
	if !v {
		t.Errorf("2020 was not a valid iteration")
		return
	}
	if c != taps.Yearly {
		t.Errorf("2020 was not yearly")
		return
	}
}

func TestGetCadenceM(t *testing.T) {
	c, v := iterCadence("2020-04 Apr")
	if !v {
		t.Errorf("2020-04 Apr was not a valid iteration")
		return
	}
	if c != taps.Monthly {
		t.Errorf("2020-04 Apr was not monthly")
		return
	}
}

func TestGetNextIterY(t *testing.T) {
	next := iterNext("2020")
	if next != "2021" {
		t.Fatalf("expected 2021; got %v", next)
	}
}

func TestGetNextIterQ(t *testing.T) {
	next := iterNext("2017 Q3")
	if next != "2017 Q4" {
		t.Fatalf("expected 2017 Q4; got %v", next)
	}
}

func TestGetNextIterM(t *testing.T) {
	next := iterNext("2020-12 Dec")
	if next != "2021-01 Jan" {
		t.Errorf("expected 2021-01 Jan; got %v", next)
		return
	}
}
