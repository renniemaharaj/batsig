package cli

import "testing"

func TestSelectExactThreshold(t *testing.T) {
	cases := []struct {
		name       string
		thresholds []int
		current    int
		want       int
	}{
		{"no thresholds", []int{}, 20, -1},
		{"exact match", []int{30}, 30, 30},
		{"multiple values", []int{10, 15, 30, 40}, 30, 30},
		{"no exact match", []int{10, 15, 30, 40}, 20, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := selectExactThreshold(tc.thresholds, tc.current)
			if got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
