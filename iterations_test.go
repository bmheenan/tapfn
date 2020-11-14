package tapfn

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestIterationQToM(t *testing.T) {
	i, err := iterResulting("2020 Q3", taps.Monthly)
	if err != nil {
		t.Errorf("Got error: %v", err)
		return
	}
	if i != "2020-09 Sep" {
		t.Errorf("Expected 2020-09 Sep; got %v", i)
		return
	}
}

func TestIterationMToQ(t *testing.T) {
	i, err := iterResulting("2020-02 Feb", taps.Quarterly)
	if err != nil {
		t.Errorf("Got error: %v", err)
		return
	}
	if i != "2020 Q1" {
		t.Errorf("Expected 2020 Q1; got %v", i)
		return
	}
}

func TestIterationMToY(t *testing.T) {
	i, err := iterResulting("2021-06 Jun", taps.Yearly)
	if err != nil {
		t.Errorf("Got error: %v", err)
		return
	}
	if i != "2021" {
		t.Errorf("Expected 2021; got %v", i)
		return
	}
}

func TestIterationYToM(t *testing.T) {
	i, err := iterResulting("2019", taps.Monthly)
	if err != nil {
		t.Errorf("Got error: %v", err)
		return
	}
	if i != "2019-12 Dec" {
		t.Errorf("Expected 2019-12 Dec; got %v", i)
		return
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
	next, err := iterNext("2020")
	if err != nil {
		t.Errorf("getNextIter returned an error: %v", err)
		return
	}
	if next != "2021" {
		t.Errorf("expected 2021; got %v", next)
		return
	}
}

func TestGetNextIterQ(t *testing.T) {
	next, err := iterNext("2017 Q3")
	if err != nil {
		t.Errorf("getNextIter returned an error: %v", err)
		return
	}
	if next != "2017 Q4" {
		t.Errorf("expected 2017 Q4; got %v", next)
		return
	}
}

func TestGetNextIterM(t *testing.T) {
	next, err := iterNext("2020-12 Dec")
	if err != nil {
		t.Errorf("getNextIter returned an error: %v", err)
		return
	}
	if next != "2021-01 Jan" {
		t.Errorf("expected 2021-01 Jan; got %v", next)
		return
	}
}
