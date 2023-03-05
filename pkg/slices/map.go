package slices

func Map[I any, O any](s []I, f func(item I) O) []O {
	out := make([]O, len(s))
	for i := range s {
		out[i] = f(s[i])
	}
	return out
}

func Filter[I any](s []I, f func(item I) bool) []I {
	out := make([]I, 0)
	for i := range s {
		if f(s[i]) {
			out = append(out, s[i])
		}
	}
	return out
}
