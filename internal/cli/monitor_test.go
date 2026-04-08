package cli

import "testing"

func TestSelectLowestAboveOrEqualThreshold(t *testing.T) {
	cases := []struct {
		name       string
		thresholds []int
		current    int
		want       int
	}{
		{"no thresholds", []int{}, 20, -1},
		{"single above", []int{30}, 20, 30},
		{"multiple above", []int{10, 15, 30, 40}, 20, 30},
		{"exact match", []int{20, 30, 40}, 20, 20},
		{"none above", []int{10, 15}, 20, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := selectLowestAboveOrEqualThreshold(tc.thresholds, tc.current)
			if got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
