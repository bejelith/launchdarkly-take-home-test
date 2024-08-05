package db

import "testing"

type value struct {
	k string
	v float64
}
type avgTest struct {
	name     string
	values   []value
	averages map[string]float64
}

func TestAverageCalculation(t *testing.T) {
	tests := []avgTest{{
		name:     "Single key",
		values:   []value{{"1", 0.2}, {"1", 0.3}},
		averages: map[string]float64{"1": 0.25},
	}, {
		name:     "Multiple key",
		values:   []value{{"1", 0.2}, {"1", 0.3}, {"2", 0}, {"2", 0.5}},
		averages: map[string]float64{"1": 0.25, "2": 0.25},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := New[string]()
			for _, value := range test.values {
				db.Add(value.k, value.v)
			}
			for k, expectedAverage := range test.averages {
				_, computedAverage := db.Get(k)
				if computedAverage != expectedAverage {
					t.Fatalf("Expected %f but %f was returned by Get()", expectedAverage, computedAverage)
				}
			}
			keys := db.GetAll()
			if len(test.averages) != len(keys) {
				t.Fatalf("Found %d keys but %d were expected", len(keys), len(test.averages))
			}
		})

	}

}
