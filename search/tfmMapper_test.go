package search

import (
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/util"
	"testing"
)

func TestCreateTFMResponse(t *testing.T) {
	//TODO test case
}

func TestRound(t *testing.T) {
	//t.Skip()
	cases := []struct {
		value, unit, expected float64
	}{
		{425.4500000000000, 100, 425.45},
		{4253.4500000000000, 100, 4253.45},
		{42.4500000000000, 100, 42.45},
		{4.254500000000000, 100, 4.25},
		{0.4254500000000000, 100, 0.43},
	}

	for _, c := range cases {
		actual := util.Round(c.value, c.unit)
		if actual != c.expected {
			t.Errorf("Round(%f, %f) = %f, want %f", c.value, c.unit, actual, c.expected)
		}
	}
}

func TestCalculateElapsedFlyingTime(t *testing.T) {
	t.Skip()
	cases := []struct {
		value    string
		expected int
	}{
		{"P0Y0M0DT3H0M0.000S", 180},
		{"P0Y0M1DT0H10M0.000S", 1450},
		{"P0Y0M0DT3H5M0.000S", 185},
	}

	for _, c := range cases {
		actual := calculateElapsedFlyingTime(c.value)
		if actual != c.expected {
			t.Errorf("CalculateElapsedFlyingTime(%v) = %b, want %b", c.value, actual, c.expected)
		}
	}
}
