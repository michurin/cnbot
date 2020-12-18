package helpers

import (
	"sort"
	"strings"
)

type AccessCtl struct {
	allowed    map[int64]struct{}
	disallowed map[int64]struct{}
}

func NewAccessCtl(a, d map[int64]struct{}) AccessCtl {
	return AccessCtl{
		allowed:    a,
		disallowed: d,
	}
}

func (a AccessCtl) IsAllowed(u int64) bool {
	if a.allowed != nil {
		_, ok := a.allowed[u]
		return ok
	}
	if a.disallowed != nil {
		_, ok := a.disallowed[u]
		return !ok
	}
	return false
}

func (a AccessCtl) String() string {
	if a.allowed != nil {
		if len(a.allowed) == 0 {
			return "white list is empty (nobody can use this bot)"
		}
		return concat(a.allowed)
	}
	if a.disallowed != nil {
		if len(a.disallowed) == 0 {
			return "black list is empty (EVERYBODY can use this bot)"
		}
		return "all except " + concat(a.disallowed)
	}
	return "nobody allowed"
}

func concat(a map[int64]struct{}) string {
	v := make([]int64, len(a))
	i := 0
	for k := range a {
		v[i] = k
		i++
	}
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	s := make([]string, len(a))
	for i, t := range v {
		s[i] = Itoa(t)
	}
	return strings.Join(s, ", ")
}
