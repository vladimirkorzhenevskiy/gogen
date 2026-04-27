package slicesutil

func Search[S ~[]E, E any](s S, f func(E) bool) (E, bool) {
	for i := range s {
		if f(s[i]) {
			return s[i], true
		}
	}

	return *new(E), false
}

func Map[S ~[]E, E, R any](s S, f func(E) R) []R {
	if s == nil {
		return nil
	}

	r := make([]R, len(s))

	for i := range s {
		r[i] = f(s[i])
	}

	return r
}

func TryMap[S ~[]E, E, R any](s S, f func(E) (R, error)) ([]R, error) {
	if len(s) == 0 {
		return []R{}, nil
	}

	res := make([]R, len(s))

	for i := range s {
		r, err := f(s[i])
		if err != nil {
			return nil, err
		}

		res[i] = r
	}

	return res, nil
}

func TryEach[S ~[]E, E any](s S, f func(E) error) error {
	for i := range s {
		if err := f(s[i]); err != nil {
			return err
		}
	}

	return nil
}

func Filter[S ~[]E, E any](s S, f func(E) bool) S {
	if s == nil {
		return nil
	}

	res := make(S, 0, len(s))

	for i := range s {
		if f(s[i]) {
			res = append(res, s[i])
		}
	}

	return res
}
