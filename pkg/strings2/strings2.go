package strings2

// Unique returns a new array containing the unique values in
// vs. The original order is preserved.
func Unique(xs []string) []string {
	if len(xs) <= 100 {
		// The O(N^2) algorithm is faster than a map when N is
		// small, due to cache locality and fewer mallocs.
		ys := make([]string, 0, len(xs))
	outer:
		for _, x := range xs {
			for _, y := range ys {
				if x == y {
					continue outer
				}
			}
			ys = append(ys, x)
		}
		return ys
	}
	// N is large enough to justify the map.
	m := make(map[string]int)
	i := 0
	for _, x := range xs {
		if _, ok := m[x]; !ok {
			m[x] = i
			i++
		}
	}
	ys := make([]string, len(m))
	for x, i := range m {
		ys[i] = x
	}
	return ys
}

// SliceContains reports whether s is present in a.
func SliceContains(a []string, s string) bool {
	return Index(a, s) >= 0
}

// Index returns the index of the first instance of s in a, or -1 if s
// is not present in a.
func Index(a []string, s string) int {
	for i, t := range a {
		if s == t {
			return i
		}
	}
	return -1
}
