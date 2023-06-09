package util

func MapToSlice(dict map[string]struct{}) []string {
	l := make([]string, 0)
	for k := range dict {
		l = append(l, k)
	}

	return l
}
