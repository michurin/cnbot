package processor

func intInSlice(x int64, a []int64) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}
